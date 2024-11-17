package entity

type Organization struct {
	ID   uint64  `json:"id" db:"org_id"`
	Info OrgInfo `json:"info"`
}

type OrgInfo struct {
	Name      string  `json:"name" validate:"min=3,max=100"`
	Address   string  `json:"address" validate:"required"`
	Long      float64 `json:"long" validate:"required,longitude"`
	Lat       float64 `json:"lat" validate:"required,latitude"`
	Type      string  `json:"type" validate:"required"`
	Telephone string  `json:"telephone,omitempty" validate:"e164"`
	Social    string  `json:"social,omitempty" validate:"url"`
	About     string  `json:"about,omitempty" validate:"max=1000"`
}

// type OrgAddInfo struct {
// 	Telephone string `json:"telephone,omitempty" validate:"e164"`
// 	Social    string `json:"social,omitempty" validate:"url"`
// 	About     string `json:"about,omitempty" validate:"max=1000"`
// }

// type City struct {
// 	ID   uint64 `json:"id"`
// 	Name string `json:"name" validate:"required,min=2,max=100"`
// }
