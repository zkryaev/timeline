package entity

// здесь перечислены общие для приложения структуры

type HashCreds struct {
	Email      string
	PasswdHash string
}

type TokenData struct {
	ID    int  `json:"id"`     // ID пользователя или организации
	IsOrg bool `json:"is_org"` // Является ли это организациями (true - организация, false - пользователь)
}

type Coordinates struct {
	Lat  float64 `json:"lat,omitempty" validate:"required,latitude"`
	Long float64 `json:"long,omitempty" validate:"required,longitude"`
}
