package orgdto

import "timeline/internal/entity"

type Timetable struct {
	OrgID     int                 `json:"org_id"`
	Timetable []*entity.OpenHours `json:"timetable" validate:"required"`
}

type TimetableReq struct {
	OrgID int
	TData entity.TokenData
}
