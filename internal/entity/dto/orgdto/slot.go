package orgdto

import "timeline/internal/entity"

type Slot struct {
	WorkerScheduleID int    `json:"worker_schedule_id,omitempty"`
	WorkerID         int    `json:"worker_id,omitempty"`
	Date             string `json:"date,omitempty" validate:"date"`
	Begin            string `json:"begin,omitempty" validate:"time"`
	End              string `json:"end,omitempty" validate:"time"`
	Busy             bool   `json:"busy,omitempty"`
}

type SlotReq struct {
	WorkerID int
	OrgID    int
	TData    entity.TokenData
}

type SlotResp struct {
	SlotID int `json:"slot_id"`
	Slot
}
