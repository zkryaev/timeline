package model

// здесь перечислены общие для приложения структуры

type Credentials struct {
	Login      string
	PasswdHash string
}

type TokenMetadata struct {
	ID    uint64
	IsOrg bool
}
