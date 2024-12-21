package orgdto

type Slot struct {
	WorkerScheduleID int    `json:"worker_schedule_id,omitempty"`
	WorkerID         int    `json:"worker_id,omitempty"`
	Date             string `json:"date,omitempty" validate:"date"`
	Begin            string `json:"begin,omitempty" validate:"time"`
	End              string `json:"end,omitempty" validate:"time"`
	Busy             bool   `json:"busy,omitempty"`
}

type SlotReq struct {
	SlotID           int `json:"slot_id"`
	WorkerID         int `json:"worker_id" validate:"required"`
	OrgID            int `json:"org_id" validate:"required"`
	WorkerScheduleID int `json:"worker_schedule_id"`
}

type SlotResp struct {
	SlotID int `json:"slot_id"`
	Slot
}

type SlotUpdate struct {
	SlotID   int `json:"slot_id" validate:"required"`
	WorkerID int `json:"worker_id" validate:"required"`
	//WorkerScheduleID int  `json:"worker_schedule_id"`
	Busy bool `json:"busy" validate:"required"`
}
