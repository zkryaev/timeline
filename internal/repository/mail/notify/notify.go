package notify

import (
	"fmt"
	"timeline/internal/config"

	"gopkg.in/gomail.v2"
)

const (
	emailFont    = "Arial, sans-serif"
	textColor    = "#333"
	codeFontSize = "24px"
)

var (
	verificationEmailTemplate = fmt.Sprintf(`
	  <div style="font-family: %s; color: %s;">
		  <p>Ваш код подтверждения:</p>
		  <div style="display: inline-block; padding: 10px; border: 1px solid #ddd; border-radius: 5px; background-color: #f0f0f0; cursor: pointer;" title="Скопируйте этот код">
			  <span style="font-size: %s; font-weight: bold; color: %s;">%%s</span>
		  </div>
		  <p style="font-weight: bold;">Никому не сообщайте этот код.</p>
		  <p style="color: #777;">Вы получили это письмо, поскольку ваш адрес был указан при регистрации в сервисе Timeline.</p>
	  </div>`, emailFont, textColor, codeFontSize, textColor)
)

type Mail interface {
	SendVerifyCode(email string, code string) error
}

type MailServer struct {
	conn *gomail.Dialer
}

func New(cfg config.Mail) *MailServer {
	return &MailServer{
		conn: gomail.NewDialer("smtp."+cfg.Host, cfg.Port, cfg.User, cfg.Password),
	}
}

// Отправка на указанную почту сгенерированного кода верификации
func (s *MailServer) SendVerifyCode(email string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@timeline.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Код подтверждения для аккаунта Timeline!")

	body := fmt.Sprintf(verificationEmailTemplate, code)
	m.SetBody("text/html", body)

	if err := s.conn.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send verification code: %w", err)
	}
	return nil
}
