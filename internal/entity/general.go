package entity

// здесь перечислены общие для приложения структуры

type Credentials struct {
	Login      string `json:"login"`       // Логин пользователя
	PasswdHash string `json:"passwd_hash"` // Хеш пароля пользователя
}

type TokenMetadata struct {
	ID    uint64 `json:"id"`     // ID пользователя или организации
	IsOrg bool   `json:"is_org"` // Является ли это организациями (true - организация, false - пользователь)
}
