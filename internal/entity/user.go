package entity

type User struct {
	UserID    int    `json:"id,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	FirstName string `json:"first_name,omitempty" validate:"required,min=3,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"required,min=3,max=100"`
	Telephone string `json:"telephone,omitempty" validate:"e164"`
	City      string `json:"city,omitempty" validate:"required"`
	About     string `json:"about,omitempty" validate:"max=500"`
}
