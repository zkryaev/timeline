package mail

import (
	"context"
	"fmt"
	"sync"
	"time"
	"timeline/internal/config"
	"timeline/internal/infrastructure"

	"timeline/internal/infrastructure/models"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var (
	// Получено экспериментальным путем.
	// Верхняя замеренная граница отправки 10 секунд, в случае сервиса gmail.com
	DefaultSendTimeout  = 10 * time.Second
	DefaulMsgBucketSize = 10000
	DefaulWorkers       = 4
)

type MailServer struct {
	WriteTimeout time.Duration
	Logger       *zap.Logger
	workers      int
	msgs         chan *gomail.Message
	wg           sync.WaitGroup
	conn         *gomail.Dialer
	contexts     []models.WorkerContext
}

// При передаче 0 в параметры будут выставлены default значения
func New(cfg config.Mail, logger *zap.Logger, WriteTimeout time.Duration, MsgBucketSize, Workers int) infrastructure.Mail {
	if WriteTimeout == 0 {
		WriteTimeout = DefaultSendTimeout
	}
	if MsgBucketSize == 0 {
		MsgBucketSize = DefaulMsgBucketSize
	}
	if Workers == 0 {
		Workers = DefaulWorkers
	}
	return &MailServer{
		WriteTimeout: WriteTimeout,
		conn:         gomail.NewDialer("smtp."+cfg.Host, cfg.Port, cfg.User, cfg.Password),
		msgs:         make(chan *gomail.Message, MsgBucketSize),
		workers:      Workers,
		contexts:     make([]models.WorkerContext, Workers),
		wg:           sync.WaitGroup{},
		Logger:       logger,
	}
}

func (s *MailServer) Start() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		s.contexts[i].Context, s.contexts[i].Close = context.WithCancel(context.Background())
		go s.worker(s.contexts[i].Context, i)
	}
}

// Gracefull
func (s *MailServer) Shutdown() {
	if len(s.contexts) != 0 {
		for i := range s.contexts {
			s.contexts[i].Close()
		}
	}
	close(s.msgs)
	s.wg.Wait()
}

func (s *MailServer) SendMsg(msg *models.Message) error {
	m, err := letterAssembly(msg)
	if err != nil {
		fmt.Println("letterAssembly:", err.Error())
		return err
	}
	fmt.Println("sendmsg", m)
	select {
	case s.msgs <- m:
		// Письмо успешно помещено в очередь
		return nil
	default:
		// Очередь заполнена, возвращаем ошибку
		time.Sleep(250 * time.Millisecond)
		return fmt.Errorf("message queue is full")
	}
}

func (s *MailServer) worker(ctx context.Context, workerID int) {
	if len(s.contexts) == 0 {
		s.Logger.Error("mail server didn't launch")
	}
	DialRetries := 2
	DialFailed := 0
	DialRetryInterval := 500 * time.Millisecond
	DialTimeout := 1 * time.Second
	SendRetries := 2
	SendRetryInterval := 1 * time.Second
	defer s.wg.Done()
	var conn gomail.SendCloser
	var err error
	var open bool
	for {
		select {
		case msg := <-s.msgs:
			fmt.Println("worker", msg)
			if msg == nil {
				continue
			}
			if conn == nil {
				conn, err = s.conn.Dial()
				if err != nil {
					DialFailed++
					s.Logger.Error(
						"failed to dial smtp connection",
						zap.Error(err),
						zap.Int("workerID", workerID),
					)
					if DialFailed < DialRetries {
						time.Sleep(DialRetryInterval)
						continue
					} else {
						s.Logger.Warn("smtp dial failed, shutting down worker", zap.Int("workerID", workerID))
						s.contexts[workerID].Close()
						return
					}
				}
				open = true
			}
			for retry := 0; retry < SendRetries; retry++ {
				if err = gomail.Send(conn, msg); err != nil {
					time.Sleep(SendRetryInterval)
				} else {
					s.Logger.Info("Email successfuly sent")
					break
				}
			}
			if err != nil {
				s.Logger.Error(
					"failed to send email",
					zap.Error(err),
					zap.Int("workerID", workerID),
				)
			}
		// Close the connection to the SMTP server if no email was sent
		case <-time.After(DialTimeout):
			if open {
				if err := conn.Close(); err != nil {
					panic(err)
				}
				open = false
			}
		case <-ctx.Done():
			if conn != nil && open {
				if err := conn.Close(); err != nil {
					s.Logger.Error(
						"failed to close smtp connection",
						zap.Error(err),
						zap.Int("workerID", workerID),
					)
				}
			}
			open = false
			return
		}
	}
}
