package orgmodel

import (
	"database/sql"
	"timeline/internal/infrastructure/models"
)

type Coordinates struct {
	Lat  float64 `db:"lat"`
	Long float64 `db:"long"`
}

type OpenHours struct {
	Weekday    sql.NullInt32 `db:"weekday"`
	Open       sql.NullTime  `db:"open"`
	Close      sql.NullTime  `db:"close"`
	BreakStart sql.NullTime  `db:"break_start"`
	BreakEnd   sql.NullTime  `db:"break_end"`
}

type OrgInfo struct {
	UUID      string  `db:"uuid"`
	OrgID     int     `db:"org_id"`
	Email     string  `db:"email"`
	Name      string  `db:"name"`
	Rating    float64 `db:"rating"`
	Type      string  `db:"type"`
	City      string  `db:"city"`
	Address   string  `db:"address"`
	Telephone string  `db:"telephone"`
	About     string  `db:"about"`
	Coordinates
}

type Organization struct {
	ShowcasesURL []*models.ImageMeta
	OrgInfo
	Timetable []*OpenHours
}

type OrgsBySearch struct {
	OrgID   int     `db:"org_id"`
	Name    string  `db:"name"`
	Rating  float64 `db:"rating"`
	Type    string  `db:"type"`
	Address string  `db:"address"`
	OpenHours
	Coordinates
}

type OrgByArea struct {
	OrgID  int     `db:"org_id"`
	Name   string  `db:"name"`
	Rating float64 `db:"rating"`
	Type   string  `db:"type"`
	OpenHours
	Coordinates
}

type OrgRegister struct {
	UUID string `db:"uuid"`
	models.HashCreds
	OrgInfo
}

type SearchParams struct {
	Page   int
	Limit  int
	Offset int
	Name   string
	Type   string
	IsRateSort bool
	IsNameSort bool
}

type AreaParams struct {
	Left  Coordinates
	Right Coordinates
}
