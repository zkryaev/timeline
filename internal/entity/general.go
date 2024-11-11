package entity

// здесь перечислены общие для приложения структуры

// Credentials представляет учетные данные для авторизации
// @Description Структура, содержащая логин и пароль пользователя
type Credentials struct {
	Login      string `json:"login"`       // Логин пользователя
	PasswdHash string `json:"passwd_hash"` // Хеш пароля пользователя
}

// TokenMetadata представляет метаданные токена для пользователя или организации
// @Description Структура для метаданных, связанных с токеном, включая ID и тип (пользователь или организация)
type TokenMetadata struct {
	ID    uint64 `json:"id"`     // ID пользователя или организации
	IsOrg bool   `json:"is_org"` // Является ли это организациями (true - организация, false - пользователь)
}
