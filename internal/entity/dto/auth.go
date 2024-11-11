package dto

import "timeline/internal/entity"

type SendCodeReq struct {
	ID    int    `json:"id" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	IsOrg bool   `json:"is_org"`
}

type VerifyCodeReq struct {
	ID    int    `json:"id" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,min=3"`
	IsOrg bool   `json:"is_org"`
}

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=64"`
}

type LoginReq struct {
	Credentials
	IsOrg bool `json:"is_org"`
}

type UserRegisterReq struct {
	Credentials
	entity.UserInfo
}

type OrgRegisterReq struct {
	City string `json:"city" validate:"required,min=2,max=100"`
	Credentials
	entity.OrgInfo
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
