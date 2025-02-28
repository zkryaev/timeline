package authdto

import "timeline/internal/entity"

type SendCodeReq struct {
	ID    int    `json:"id" validate:"required,gt=0"`
	Email string `json:"email" validate:"required,email"`
	IsOrg bool   `json:"is_org"`
}

type VerifyCodeReq struct {
	ID    int    `json:"id" validate:"required,gt=0"`
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,min=3"`
	IsOrg bool   `json:"is_org"`
}

// Credentials структура для хранения данных авторизации
type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=64"`
}

type LoginReq struct {
	Credentials
	IsOrg bool `json:"is_org"`
}

type UserRegisterReq struct {
	UUID string
	Credentials
	entity.User
}

type RegisterResp struct {
	ID int `json:"id"`
}

type OrgRegisterReq struct {
	UUID string
	Credentials
	entity.OrgInfo
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	Token string `json:"access_token"`
}
