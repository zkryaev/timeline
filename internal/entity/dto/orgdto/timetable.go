package orgdto

import "timeline/internal/entity"

type Timetable struct {
	OrgID     int                 `json:"org_id" validate:"required"`
	Timetable []*entity.OpenHours `json:"timetable" validate:"required"`
}
