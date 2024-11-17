package models

type UserRegisterModel struct {
	HashCreds
	UserInfo
}

type UserInfo struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Telephone string `db:"telephone"`
	Social    string `db:"social"`
	About     string `db:"about"`
}
