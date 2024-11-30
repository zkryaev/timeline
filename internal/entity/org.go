package entity

type OrgInfo struct {
	Name      string    `json:"name" validate:"min=3,max=100"`
	Rating    float64   `json:"rating,omitempty"`
	Address   string    `json:"address" validate:"required"`
	Long      float64   `json:"long" validate:"required,longitude"`
	Lat       float64   `json:"lat" validate:"required,latitude"`
	Type      string    `json:"type" validate:"required"`
	Telephone string    `json:"telephone" validate:"e164"`
	City      string    `json:"city" validate:"required"`
	About     string    `json:"about,omitempty" validate:"max=1500"`
	Timetable *OpenHours `json:"timetable,omitempty"`
}

type MapOrgInfo struct {
	OrgID         int        `json:"org_id"`
	Name          string     `json:"name"`
	Rating        float64    `json:"rating"`
	Type          string     `json:"type"`
	TodaySchedule *OpenHours `json:"today_schedule"`
	Coords        MapPoint   `json:"coords"`
}

type OpenHours struct {
	Weekday    int    `json:"weekday"`
	Open       string `json:"open,omitempty" validate:"time"`
	Close      string `json:"close,omitempty" validate:"time"`
	BreakStart string `json:"break_start,omitempty" validate:"time"`
	BreakEnd   string `json:"break_end,omitempty" validate:"time"`
}

// type OrgAddInfo struct {
// 	Telephone string `json:"telephone,omitempty" validate:"e164"`
// 	Social    string `json:"social,omitempty" validate:"url"`
// 	About     string `json:"about,omitempty" validate:"max=1000"`
// }

// type City struct {
// 	ID   uint64 `json:"id"`
// 	Name string `json:"name" validate:"required,min=2,max=100"`
// }
