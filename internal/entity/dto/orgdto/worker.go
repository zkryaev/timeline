package orgdto

import "timeline/internal/entity"

type AddWorkerReq struct {
	OrgID      int           `json:"org_id" validate:"required"`
	WorkerInfo entity.Worker `json:"worker_info" validate:"required"`
}

type UpdateWorkerReq struct {
	WorkerID   int           `json:"worker_id" validate:"required"`
	OrgID      int           `json:"org_id" validate:"required"`
	WorkerInfo entity.Worker `json:"worker_info"`
}

type AssignWorkerReq struct {
	ServiceID int `json:"service_id" validate:"required"`
	WorkerID  int `json:"worker_id" validate:"required"`
	OrgID     int `json:"org_id" validate:"required"`
}

type WorkerResp struct {
	WorkerID   int            `json:"worker_id"`
	OrgID      int            `json:"org_id,omitempty"`
	WorkerInfo *entity.Worker `json:"worker_info,omitempty"`
}

type WorkerList struct {
	List  []*WorkerResp `json:"worker_list"`
	Found int           `json:"found"`
}
