package models

type UserRegister struct {
	HashCreds
	UserInfo
}

type UserInfo struct {
	OrgID     int    `db:"user_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Telephone string `db:"telephone"`
	City      string `db:"city"`
	About     string `db:"about"`
}
