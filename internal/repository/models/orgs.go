package models

type OrgRegisterModel struct {
	HashCreds
	OrgInfo
}

type OrgInfo struct {
	Name      string  `db:"name"`
	Address   string  `db:"org_address"`
	Long      float64 `db:"long"`
	Lat       float64 `db:"lat"`
	Type      string  `db:"type"`
	Telephone string  `db:"telephone"`
	Social    string  `db:"social"`
	About     string  `db:"about"`
}

type VerifyInfo struct {
	userID string `db:"user_id"`
	code   string `db:"code"`
}
