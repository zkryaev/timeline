package usermodel

import "timeline/internal/repository/models"

type UserRegister struct {
	models.HashCreds
	UserInfo
}

type UserInfo struct {
	UserID    int    `db:"user_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Telephone string `db:"telephone"`
	City      string `db:"city"`
	About     string `db:"about"`
}
