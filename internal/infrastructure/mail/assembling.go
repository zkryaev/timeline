package mail

import (
	"errors"
	"fmt"
	"io"
	"time"
	"timeline/internal/infrastructure/models"

	ics "github.com/arran4/golang-ical"
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
	 	  <span style="font-weight: bold; color: %s">Дата и время записи:</span>
	  	  <span>%%s</span>
		</div>
		  <p>Ждем вас!</p>
	  </div>
	  `, emailFont, textColor, labelColor, labelColor, labelColor)

	cancellationTemplate = fmt.Sprintf(`
	  <div style="font-family: %s; color: %s; line-height: 1.6; margin: 20px; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd;">
		  <p>Здравствуйте!</p>
		  <p>Ваша запись на услугу была отменена.</p>
		  <div style="margin-bottom: 5px;">
			  <span style="font-weight: bold; color: %s">Организация:</span>
			  <span>%%s</span>
		  </div>
		  <div style="margin-bottom: 5px;">
			  <span style="font-weight: bold; color: %s">Услуга:</span>
			  <span>%%s</span>
		  </div>
		  <div style="margin-bottom: 5px;">
			  <span style="font-weight: bold; color: %s">Дата и время записи:</span>
			  <span>%%s</span>
		  </div>
		  <div style="margin-bottom: 5px;">
			<span style="font-weight: bold; color: %s">Причина отмены:</span>
			<span>%%s</span>
		</div>
		  <p>Если отмена произошла по ошибке или у вас есть вопросы, пожалуйста, свяжитесь с организацией.</p>
		  <p style="color: #777; font-size: 0.9em;">Это автоматическое уведомление. Пожалуйста, не отвечайте на это письмо.</p>
	  </div>`,
		emailFont, textColor, labelColor, labelColor, labelColor)
)

var (
	ErrWrongMsgType = errors.New("set wrong msg type")
	ErrTypeNotExist = errors.New("set type not exist")
)

var (
	VerificationType = "verification"
	ReminderType     = "reminder"
	CancelationType  = "cancelation"
)

// Сборка письма
func letterAssembly(data *models.Message) (*gomail.Message, error) {
	var body, subject string
	var icsContent string
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
		fields, ok := data.Value.(*models.ReminderMsg)
		if !ok {
			return nil, ErrWrongMsgType
		}
		body = fmt.Sprintf(reminderTemplate,
			fields.Organization,
			fields.Service,
			fields.SessionDate.Format("02.01.2006")+" : "+fields.SessionStart.Format("15:04")+"-"+fields.SessionEnd.Format("15:04"),
		)
		if data.IsAttach {
			icsContent = icsCreate(fields)
		}
	case CancelationType:
		subject = "Ваша запись отменена"
		fields, ok := data.Value.(*models.CancelMsg)
		if !ok {
			return nil, ErrWrongMsgType
		}
		body = fmt.Sprintf(reminderTemplate,
			fields.Organization,
			fields.Service,
			fields.SessionDate+" : "+fields.SessionStart+"-"+fields.SessionEnd,
			fields.CancelReason,
		)
	default:
		return nil, ErrTypeNotExist
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@timeline.ru")
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if data.IsAttach {
		m.Attach("event.ics", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write([]byte(icsContent))
			return err
		}))
	}
	return m, nil
}

func icsCreate(msg *models.ReminderMsg) string {
	eventUID := fmt.Sprintf("%s@%s", time.Now().Format("20060102T150405Z"), "timeline.ru")
	reminderUID := fmt.Sprintf("reminder-%s", eventUID)

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	event := cal.AddEvent(eventUID)
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(msg.SessionStart)
	event.SetEndAt(msg.SessionEnd)
	event.SetSummary(msg.Service)
	event.SetLocation(msg.Address)
	event.SetDescription(msg.ServiceDesc)
	event.SetOrganizer("timeline@gmail.com", ics.WithCN(msg.Organization))
	reminder := ics.NewAlarm(reminderUID)
	reminder.SetSummary(msg.Service)
	reminder.SetDescription("Ваша запись начинается через 2 часа!")
	reminder.SetAction(ics.ActionDisplay)
	reminder.SetTrigger("-PT2H") // за 2 часс до события напомнить
	event.AddVAlarm(reminder)

	return cal.Serialize()
}
