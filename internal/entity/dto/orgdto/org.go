package orgdto

import "timeline/internal/entity"

type Organization struct {
	OrgID     int                 `json:"id"`
	Info      *entity.OrgInfo     `json:"info"`
	Timetable []*entity.OpenHours `json:"timetable,omitempty"`
}

type OrgUpdateReq struct {
	OrgID int `json:"org_id" validate:"required"`
	entity.OrgInfo
	Timetable []*entity.OpenHours `json:"timetable,omitempty"`
}

type TimetableUpdate struct {
	OrgID     int                 `json:"org_id" validate:"required"`
	Timetable []*entity.OpenHours `json:"timetable" validate:"required"`
}
