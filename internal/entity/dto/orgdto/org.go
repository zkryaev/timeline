package orgdto

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/s3dto"
)

type Organization struct {
	OrgID     int `json:"id"`
	ImagesURL []*s3dto.FileURL
	Info      *entity.OrgInfo     `json:"info"`
	Timetable []*entity.OpenHours `json:"timetable,omitempty"`
}

type OrgUpdateReq struct {
	OrgID int `json:"org_id" validate:"required"`
	entity.OrgInfo
	Timetable []*entity.OpenHours `json:"timetable,omitempty"`
}
