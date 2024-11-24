package models

type OrgRegister struct {
	HashCreds
	OrgInfo
}
type OrgInfo struct {
	OrgID     int     `db:"org_id"`
	Name      string  `db:"name"`
	Rating    float64 `db:"rating"`
	Type      string  `db:"type"`
	City      string  `db:"city"`
	Address   string  `db:"address"`
	Telephone string  `db:"telephone"`
	Lat       float64 `db:"lat"`
	Long      float64 `db:"long"`
	About     string  `db:"about"`
}

type OrgSummary struct {
	OrgID  int     `db:"org_id"`
	Name   string  `db:"name"`
	Rating float64 `db:"rating"`
	Type   string  `db:"type"`
}

type OrgUpdate struct {
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