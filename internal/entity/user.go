package entity

type User struct {
	ID   uint64   `json:"id" db:"user_id"`
	Info UserInfo `json:"info"`
}

type UserInfo struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=100"`
	LastName  string `json:"last_name" validate:"required,min=3,max=100"`
	Telephone string `json:"telephone" validate:"e164"`
	Social    string `json:"social" validate:"url"`
	About     string `json:"about" validate:"max=500"`
}
