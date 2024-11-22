package models

type OrgRegister struct {
	HashCreds
	OrgInfo
}
type OrgInfo struct {
	OrgID     int     `db:"org_id"`
	Name      string  `db:"name"`
	Type      string  `db:"type"`
	City      string  `db:"city"`
	Address   string  `db:"address"`
	Telephone string  `db:"telephone"`
	Lat       float64 `db:"lat"`
	Long      float64 `db:"long"`
	About     string  `db:"about"`
}

type VerifyInfo struct {
	userID string `db:"user_id"`
	code   string `db:"code"`
}
