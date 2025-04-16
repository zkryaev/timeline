package userdto

type UserUpdateReq struct {
	UserID    int
	FirstName string `json:"first_name" validate:"min=3,max=100"`
	LastName  string `json:"last_name" validate:"min=3,max=100"`
	Telephone string `json:"telephone" validate:"e164"`
	City      string `json:"city"`
	About     string `json:"about" validate:"max=500"`
}
