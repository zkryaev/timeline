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

type RecordParam struct {
	RecordID int
	TData    entity.TokenData
}

type RecordListParams struct {
	OrgID    int  `json:"org_id"`
	UserID   int  `json:"user_id"`
	Fresh    bool `json:"fresh"`
	Reviewed bool `json:"reviewed"`
	Limit    int  `json:"limit"`
	Page     int  `json:"page"`
	TData    entity.TokenData
}

type RecordScrap struct {
	RecordID int                  `json:"record_id"`
	Reviewed bool                 `json:"reviewed"`
	Org      *orgdto.Organization `json:"org,omitempty"`
	User     *entity.User         `json:"user,omitempty"`
	Slot     *orgdto.Slot         `json:"slot,omitempty"`
	Service  *entity.Service      `json:"service,omitempty"`
	Worker   *entity.Worker       `json:"worker,omitempty"`
	Feedback *Feedback            `json:"feedback,omitempty"`
}

type RecordList struct {
	List  []*RecordScrap `json:"record_list"`
	Found int            `json:"found"`
}

type RecordCancelation struct {
	TData        entity.TokenData
	UserID       int    `json:"user_id"`
	OrgID        int    `json:"org_id"`
	RecordID     int    `json:"record_id"`
	CancelReason string `json:"cancel_reason"`
}
