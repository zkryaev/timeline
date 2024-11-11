package entity

type User struct {
	ID   uint64   `json:"id"`   // Уникальный идентификатор пользователя
	Info UserInfo `json:"info"` // Информация о пользователе
}

type UserInfo struct {
	Name      string `json:"name" validate:"required,min=3,max=100"` // Имя пользователя
	Telephone string `json:"telephone" validate:"required,e164"`     // Телефон пользователя
	Social    string `json:"social" validate:"required,url"`         // Социальная ссылка
	About     string `json:"about" validate:"max=500"`               // Описание пользователя
}
