package entity


type Organization struct {
	ID   uint64  `json:"id"`   // Уникальный идентификатор организации
	Info OrgInfo `json:"info"` // Информация о организации
}


type OrgInfo struct {
	Name      string  `json:"name" validate:"required,min=3,max=100"` // Название организации
	Telephone string  `json:"telephone" validate:"required,e164"`     // Телефон организации
	Social    string  `json:"social" validate:"url"`                  // Социальная ссылка
	About     string  `json:"about" validate:"max=1000"`              // Описание организации
	Address   string  `json:"address" validate:"required"`            // Адрес организации
	Long      float64 `json:"long" validate:"longitude"`              // Долгота
	Lat       float64 `json:"lat" validate:"latitude"`                // Широта
}

// type City struct {
// 	ID   uint64 `json:"id"`
// 	Name string `json:"name" validate:"required,min=2,max=100"`
// }
