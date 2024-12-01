package models

import (
	"database/sql"
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
	OrgID     int     `db:"org_id"`
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
	HashCreds
	OrgInfo
}
