package usermodel

import "timeline/internal/infrastructure/models"

type UserRegister struct {
	UUID string `db:"uuid"`
	models.HashCreds
	UserInfo
}

type UserInfo struct {
	UserID    int    `db:"user_id"`
	UUID      string `db:"uuid"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Telephone string `db:"telephone"`
	City      string `db:"city"`
	About     string `db:"about"`
}
