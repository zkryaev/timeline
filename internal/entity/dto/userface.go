package dto

import "timeline/internal/entity"

type SearchReq struct {
	Page  int    `json:"page" validator:"required,gte=1"`
	Limit int    `json:"limit" validator:"required,gt=0"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type SearchResp struct {
	Orgs []*entity.Organization `json:"orgs"`
}

type OrgAreaReq struct {
	LeftLowerCorner  entity.MapPoint `json:"left_lower_corner"`
	RightUpperCorner entity.MapPoint `json:"right_upper_corner"`
}

type OrgAreaResp struct {
	Orgs []*entity.MapOrgInfo `json:"map_orgs"`
}
