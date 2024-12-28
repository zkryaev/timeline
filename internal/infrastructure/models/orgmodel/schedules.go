package orgmodel

import "time"

type WorkerSchedule struct {
	WorkerID        int `db:"worker_id"`
	OrgID           int `db:"org_id"`
	SessionDuration int `db:"session_duration"`
	Schedule        []*Schedule
}

type ScheduleList struct {
	Workers []*WorkerSchedule
	Found   int
}

type Schedule struct {
	WorkerScheduleID int       `db:"worker_schedule_id"`
	Weekday          int       `db:"weekday"`
	Start            time.Time `db:"start"`
	Over             time.Time `db:"over"`
}

type ScheduleParams struct {
	WorkerID int `db:"worker_id"`
	OrgID    int `db:"org_id"`
	Weekday  int `db:"weekday"`
	Limit    int
	Offset   int
}
