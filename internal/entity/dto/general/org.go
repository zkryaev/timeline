package general

import (
	"timeline/internal/entity"
)

type SearchReq struct {
	Page  int    `json:"page" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"required,min=1"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}

type SearchResp struct {
	Found int                    `json:"found"`
	Orgs  []*entity.OrgsBySearch `json:"orgs"`
}

type OrgAreaReq struct {
	LeftLowerCorner  entity.Coordinates `json:"left_lower_corner"`
	RightUpperCorner entity.Coordinates `json:"right_upper_corner"`
}

type OrgAreaResp struct {
	Found int                  `json:"found"`
	Orgs  []*entity.MapOrgInfo `json:"map_orgs"`
}
