package orgdto

import "timeline/internal/entity"

type OrgUpdateReq struct {
	OrgID     int     `json:"org_id" validate:"required"`
	Name      string  `json:"name" validate:"min=3,max=100"`
	Address   string  `json:"address"`
	Long      float64 `json:"long" validate:"longitude"`
	Lat       float64 `json:"lat" validate:"latitude"`
	Type      string  `json:"type"`
	Telephone string  `json:"telephone" validate:"e164"`
	City      string  `json:"city"`
	About     string  `json:"about,omitempty" validate:"max=1500"`
}

type OrgUpdateResp struct {
	OrgID     int     `json:"org_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Long      float64 `json:"long"`
	Lat       float64 `json:"lat"`
	Type      string  `json:"type"`
	Telephone string  `json:"telephone"`
	City      string  `json:"city"`
	About     string  `json:"about,omitempty"`
}

type SearchReq struct {
	Page  int    `json:"page" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"required,min=1"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
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
