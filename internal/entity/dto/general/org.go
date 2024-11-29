package general

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
)

type SearchReq struct {
	Page  int    `json:"page" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"required,min=1"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}

type SearchResp struct {
	Pages int                    `json:"pages"`
	Orgs  []*orgdto.Organization `json:"orgs"`
}

type OrgAreaReq struct {
	LeftLowerCorner  entity.MapPoint `json:"left_lower_corner"`
	RightUpperCorner entity.MapPoint `json:"right_upper_corner"`
}

type OrgAreaResp struct {
	Found int                  `json:"found"`
	Orgs  []*entity.MapOrgInfo `json:"map_orgs"`
}
