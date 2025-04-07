package orgmodel

import "time"

type Slot struct {
	SlotID           int       `db:"slot_id"`
	WorkerScheduleID int       `db:"worker_schedule_id"`
	WorkerID         int       `db:"worker_id"`
	Date             time.Time `db:"date"`
	Begin            time.Time `db:"session_begin"`
	End              time.Time `db:"session_end"`
	Busy             bool      `db:"busy"`
}

type SlotsMeta struct {
	SlotID   int `db:"slot_id"`
	WorkerID int `db:"worker_id"`
	UserID   int
	OrgID    int
	// WorkerScheduleID int `db:"worker_schedule_id"`
}
