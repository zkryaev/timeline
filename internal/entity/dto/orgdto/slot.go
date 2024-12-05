package orgdto

type SlotInfo struct {
	WorkerScheduleID int    `json:"worker_schedule_id"`
	WorkerID         int    `json:"worker_id"`
	Date             string `json:"date" validate:"date"`
	Begin            string `json:"begin" validate:"time"`
	End              string `json:"end" validate:"time"`
	Busy             bool   `json:"busy"`
}

type SlotReq struct {
	SlotID           int `json:"slot_id"`
	WorkerID         int `json:"worker_id" validate:"required"`
	WorkerScheduleID int `json:"worker_schedule_id"`
}

type SlotResp struct {
	SlotID int `json:"slot_id"`
	SlotInfo
}

type SlotUpdate struct {
	SlotID           int  `json:"slot_id" validate:"required"`
	WorkerID         int  `json:"worker_id" validate:"required"`
	WorkerScheduleID int  `json:"worker_schedule_id"`
	Busy             bool `json:"busy" validate:"required"`
}
