package notify

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"timeline/internal/config"
	"timeline/internal/repository/mail"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

const (
	emailFont    = "Arial, sans-serif"
	textColor    = "#333"
	codeFontSize = "24px"
	labelColor   = "#000"
)

var (
	verificationTemplate = fmt.Sprintf(`
	  <div style="font-family: %s; color: %s;">
		  <p>Ваш код подтверждения:</p>
		  <div style="display: inline-block; padding: 10px; border: 1px solid #ddd; border-radius: 5px; background-color: #f0f0f0; cursor: pointer;" title="Скопируйте этот код">
			  <span style="font-size: %s; font-weight: bold; color: %s;">%%s</span>
		  </div>
		  <p style="font-weight: bold;">Никому не сообщайте этот код.</p>
		  <p style="color: #777;">Вы получили это письмо, поскольку ваш адрес был указан при регистрации в сервисе Timeline.</p>
	  </div>`, emailFont, textColor, codeFontSize, textColor)

	reminderTemplate = fmt.Sprintf(`
	  <div style="font-family: %s; color: %s; line-height: 1.6; margin: 20px; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd;">
		<p>Здравствуйте!</p>
		<p>Напоминаем вам о вашей записи на услугу.</p>
		<div style="margin-bottom: 5px;">
		  <span style="font-weight: bold; color: %s">Организация:</span>
		  <span>%%s</span>
		</div>
		<div style="margin-bottom: 5px;">
		  <span style="font-weight: bold; color: %s">Услуга:</span>
		  <span>%%s</span>
		</div>
		<div style="margin-bottom: 5px;">
		  <span style="font-weight: bold; color: %s">Время:</span>
		  <span>%%s</span>
		</div>
		  <p>Ждем вас!</p>
	  </div>
	  `, emailFont, textColor, labelColor, labelColor, labelColor)
)

var (
	ErrWrongMsgType = errors.New("set wrong msg type")
	ErrTypeNotExist = errors.New("set type not exist")
)

var (
	// Получено экспериментальным путем.
	// Верхняя замеренная граница отправки 10 секунд, в случае сервиса gmail.com
	DefaultSendTimeout  = 10 * time.Second
	DefaulMsgBucketSize = 100
	DefaulWorkers       = 4
)

var (
	VerificationType = "verification"
	ReminderType     = "reminder"
)

type Mail interface {
	//SendVerifyCode(email string, code string) error
	SendMsg(msg *mail.Message) error
	Start()
	Shutdown()
}

type CtxCloser struct {
	Context context.Context
	Close   context.CancelFunc
}

type MailServer struct {
	WriteTimeout time.Duration
	Logger       *zap.Logger
	workers      int
	msgs         chan *gomail.Message
	wg           sync.WaitGroup
	conn         *gomail.Dialer
	contexts     []CtxCloser
}

// По умолчанию WriteTimeout = 10 секундам
func New(cfg config.Mail, logger *zap.Logger, WriteTimeout time.Duration, MsgBucketSize, Workers int) *MailServer {
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
		contexts:     make([]CtxCloser, Workers),
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

func (s *MailServer) worker(ctx context.Context, workerID int) {
	if len(s.contexts) == 0 {
		s.Logger.Error("mail server didn't launch")
	}
	DialRetries := 2
	DialFailed := 0
	DialRetryInterval := 500 * time.Millisecond
	defer s.wg.Done()
	var conn gomail.SendCloser
	var err error
	for {
		select {
		case msg := <-s.msgs:
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
			}
			maxRetries := 2
			retryInterval := 1 * time.Second
			for retry := 0; retry < maxRetries; retry++ {
				if err = gomail.Send(conn, msg); err != nil {
					time.Sleep(retryInterval)
				} else {
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
		case <-ctx.Done():
			if conn != nil {
				if err := conn.Close(); err != nil {
					s.Logger.Error(
						"failed to close smtp connection",
						zap.Error(err),
						zap.Int("workerID", workerID),
					)
				}
			}
			return
		}
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

// Сборка письма
func letterAssembly(data *mail.Message) (*gomail.Message, error) {
	var body, subject string
	switch data.Type {
	case VerificationType:
		subject = "Код подтверждения для аккаунта Timeline!"
		code, ok := data.Value.(string)
		if !ok {
			return nil, ErrWrongMsgType
		}
		body = fmt.Sprintf(verificationTemplate, code)
	case ReminderType:
		subject = "Напоминание о вашей записи!"
		fields, ok := data.Value.(mail.ReminderFields)
		if !ok {
			return nil, ErrWrongMsgType
		}
		body = fmt.Sprintf(reminderTemplate,
			fields.Organization,
			fields.Service,
			fields.SessionTime,
		)
	default:
		return nil, ErrTypeNotExist
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@timeline.ru")
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return m, nil
}

func (s *MailServer) SendMsg(msg *mail.Message) error {
	m, err := letterAssembly(msg)
	if err != nil {
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

// Отправка на указанную почту сгенерированного кода верификации
// func (s *MailServer) SendVerifyCode(email string, code string) error {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", "noreply@timeline.ru")
// 	m.SetHeader("To", email)
// 	m.SetHeader("Subject", "Код подтверждения для аккаунта Timeline!")

// 	body := fmt.Sprintf(verificationTemplate, code)
// 	m.SetBody("text/html", body)

// 	if err := s.conn.DialAndSend(m); err != nil {
// 		return fmt.Errorf("failed to send verification code: %w", err)
// 	}
// 	return nil
// }
