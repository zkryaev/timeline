package entity

type Organization struct {
	ID   uint64  `json:"id"`
	Info OrgInfo `json:"info"`
}

type OrgInfo struct {
	Name      string  `json:"name" validate:"required,min=3,max=100"`
	Telephone string  `json:"telephone" validate:"required,e164"`
	Social    string  `json:"social" validate:"url"`
	About     string  `json:"about" validate:"max=1000"`
	Address   string  `json:"address" validate:"required"`
	Long      float64 `json:"long" validate:"longitude"`
	Lat       float64 `json:"lat" validate:"latitude"`
}

type City struct {
	ID   uint64 `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=100"`
}
