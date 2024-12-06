package recordto

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
)

type Record struct {
	RecordID  int  `json:"record_id"`
	OrgID     int  `json:"org_id"`
	UserID    int  `json:"user_id"`
	SlotID    int  `json:"slot_id"`
	ServiceID int  `json:"service_id"`
	WorkerID  int  `json:"worker_id"`
	Reviewed  bool `json:"reviewed"`
}

type RecordListParams struct {
	OrgID    int  `json:"org_id"`
	UserID   int  `json:"user_id"`
	Reviewed bool `json:"reviewed"`
}

type RecordScrap struct {
	RecordID int             `json:"record_id"`
	Reviewed bool            `json:"reviewed"`
	Org      *entity.OrgInfo `json:"org,omitempty"`
	User     *entity.User    `json:"user,omitempty"`
	Slot     *orgdto.Slot    `json:"slot,omitempty"`
	Service  *entity.Service `json:"service,omitempty"`
	Worker   *entity.Worker  `json:"worker,omitempty"`
	Feedback *Feedback       `json:"feedback,omitempty"`
}

type RecordList struct {
	List []*RecordScrap
}
