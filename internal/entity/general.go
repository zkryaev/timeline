package entity

// здесь перечислены общие для приложения структуры

type HashCreds struct {
	Email      string
	PasswdHash string
}

type TokenMetadata struct {
	ID    uint64 `json:"id"`     // ID пользователя или организации
	IsOrg bool   `json:"is_org"` // Является ли это организациями (true - организация, false - пользователь)
}

type MapPoint struct {
	Lat  float64 `json:"lat" validate:"required,latitude"`
	Long float64 `json:"long" validate:"required,longitude"`
}
