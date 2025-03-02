package orgdto

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/s3dto"
)

type Organization struct {
	OrgID        int                 `json:"id"`
	UUID         string              `json:"uuid,omitempty"`
	ShowcasesURL []*s3dto.FileURL    `json:"showcases_url,omitempty"`
	Info         *entity.OrgInfo     `json:"info,omitempty"`
	Timetable    []*entity.OpenHours `json:"timetable,omitempty"`
}

type OrgUpdateReq struct {
	OrgID int `json:"org_id" validate:"required"`
	entity.OrgInfo
	Timetable []*entity.OpenHours `json:"timetable,omitempty"`
}
