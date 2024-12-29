package orgdto

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/s3dto"
)

type Organization struct {
	OrgID        int                 `json:"id"`
	UUID         string              `json:"uuid"`
	ShowcasesURL []*s3dto.FileURL    `json:"showcases_url"`
	Info         *entity.OrgInfo     `json:"info"`
	Timetable    []*entity.OpenHours `json:"timetable,omitempty"`
}

type OrgUpdateReq struct {
	OrgID int `json:"org_id" validate:"required"`
	entity.OrgInfo
	Timetable []*entity.OpenHours `json:"timetable,omitempty"`
}
