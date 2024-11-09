package dto

import "timeline/internal/model"

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=64,containsupper,containslower,containsnumber,containsany=@#_"`
}

type LoginReq struct {
	Credentials
	IsOrg bool `json:"is_org" validate:"required"`
}

type UserRegisterReq struct {
	Credentials
	model.UserInfo
}

type OrgRegisterReq struct {
	Credentials
	model.OrgInfo
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
