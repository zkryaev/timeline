package userdto

import "timeline/internal/entity"

type UserGetResp struct {
	UserID int `json:"user_id"`
	entity.UserInfo
}

type UserUpdateReq struct {
	UserID    int    `json:"id"`
	FirstName string `json:"first_name" validate:"min=3,max=100"`
	LastName  string `json:"last_name" validate:"min=3,max=100"`
	Telephone string `json:"telephone" validate:"e164"`
	City      string `json:"city"`
	About     string `json:"about" validate:"max=500"`
}

type UserUpdateResp struct {
	UserID    int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Telephone string `json:"telephone"`
	City      string `json:"city"`
	About     string `json:"about"`
}
