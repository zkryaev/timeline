package orgmodel

import (
	"time"
	"timeline/internal/infrastructure/models"
)

type Slot struct {
	SlotID           int       `db:"slot_id"`
	WorkerScheduleID int       `db:"worker_schedule_id"`
	WorkerID         int       `db:"worker_id"`
	Date             time.Time `db:"date"`
	Begin            time.Time `db:"session_begin"`
	End              time.Time `db:"session_end"`
	Busy             bool      `db:"busy"`
}

type SlotsReq struct {
	WorkerID int `db:"worker_id"`
	OrgID    int `db:"org_id"`
	TData    models.TokenData
}
