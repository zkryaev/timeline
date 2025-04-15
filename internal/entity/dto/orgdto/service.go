package orgdto

import "timeline/internal/entity"

type AddServiceReq struct {
	OrgID       int            `json:"org_id"`
	ServiceInfo entity.Service `json:"service_info" validate:"required"`
}

type UpdateServiceReq struct {
	OrgID       int
	ServiceID   int            `json:"service_id" validate:"required"`
	ServiceInfo entity.Service `json:"service_info" validate:"required"`
}

type ServiceResp struct {
	ServiceID   int             `json:"service_id"`
	OrgID       int             `json:"org_id,omitempty"`
	ServiceInfo *entity.Service `json:"service_info,omitempty"`
}

type ServiceList struct {
	List  []*ServiceResp `json:"service_list"`
	Found int            `json:"found"`
}
