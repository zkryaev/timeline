package models

type Service struct {
	ServiceID   int     `db:"service_id"`
	OrgID       int     `db:"org_id"`
	Name        string  `db:"name"`
	Cost        float64 `db:"cost"`
	Description string  `db:"description"`
}
