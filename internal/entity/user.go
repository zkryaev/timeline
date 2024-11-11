package entity

type User struct {
	ID   uint64   `json:"id"`
	Info UserInfo `json:"info"`
}

type UserInfo struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Telephone string `json:"telephone" validate:"required,e164"`
	Social    string `json:"social" validate:"required,url"`
	About     string `json:"about" validate:"max=500"`
}
