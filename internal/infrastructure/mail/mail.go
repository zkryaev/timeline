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

const (
	// Получено экспериментальным путем.
	// Верхняя замеренная граница отправки 10 секунд, в случае сервиса gmail.com
	DefaultSendTimeout  = 10 * time.Second
	DefaulMsgBucketSize = 10000
	DefaulWorkers       = 4

	DialAttempts = 2
	DialInterval = 500 * time.Millisecond
	DialTimeout  = 1 * time.Second

	SendAttempts = 2
	SendInterval = 1 * time.Second

	CloseInterval = 1 * time.Second
	CloseAttempts = 2
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
func New(cfg config.Mail, logger *zap.Logger, writeTimeout time.Duration, msgBucketSize, workers int) infrastructure.Mail {
	if writeTimeout == 0 {
		writeTimeout = DefaultSendTimeout
	}
	if msgBucketSize == 0 {
		msgBucketSize = DefaulMsgBucketSize
	}
	if workers == 0 {
		workers = DefaulWorkers
	}
	return &MailServer{
		WriteTimeout: writeTimeout,
		conn:         gomail.NewDialer("smtp."+cfg.Host, cfg.Port, cfg.User, cfg.Password),
		msgs:         make(chan *gomail.Message, msgBucketSize),
		workers:      workers,
		contexts:     make([]models.WorkerContext, workers),
		wg:           sync.WaitGroup{},
		Logger:       logger,
	}
}

func (s *MailServer) Start() {
	for i := range s.workers {
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
		// TODO: логировать ошибку
		return err
	}
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
	closer := func(conn gomail.SendCloser, isOpen bool) {
		if conn != nil && isOpen {
			for range CloseAttempts {
				err := conn.Close()
				if err == nil {
					break
				}
				s.Logger.Fatal("closer", zap.String("failed to close conn", err.Error()))
				time.Sleep(CloseInterval)
			}
		}
	}
	defer s.wg.Done()
	var conn gomail.SendCloser
	var err error
	var open bool
	for {
		select {
		case msg := <-s.msgs:
			if msg == nil {
				continue
			}
			if conn == nil {
				var failedDialCnt int
				for failedDialCnt < DialAttempts {
					conn, err = s.conn.Dial()
					if err == nil {
						break
					}
					s.Logger.Error(
						"failed to dial smtp connection",
						zap.Error(err),
						zap.Int("workerID", workerID),
					)
					failedDialCnt++
					time.Sleep(DialInterval)
				}
				if failedDialCnt >= DialAttempts {
					s.Logger.Warn(fmt.Sprintf("smtp dial failed, stop running worker_%d...", workerID))
					s.contexts[workerID].Close()
					return
				}
				open = true
			}
			for range SendAttempts {
				if err = gomail.Send(conn, msg); err != nil {
					time.Sleep(SendInterval)
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
		case <-time.After(DialTimeout):
			closer(conn, open)
			return
		case <-ctx.Done():
			closer(conn, open)
			return
		}
	}
}
