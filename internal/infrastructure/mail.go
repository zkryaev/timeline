package infrastructure

import "timeline/internal/infrastructure/models"

type Mail interface {
	SendMsg(msg *models.Message) error
	Start()
	Shutdown()
}
