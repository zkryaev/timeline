package models

import (
	"database/sql"
)

type Coordinate struct {
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
	Timetable []*OpenHours
}

type OrgSummary struct {
	OrgID  int     `db:"org_id"`
	Name   string  `db:"name"`
	Rating float64 `db:"rating"`
	Type   string  `db:"type"`
	OpenHours
	Coordinate
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
	Timetable []*OpenHours
}

type OrgRegister struct {
	HashCreds
	OrgInfo
}
