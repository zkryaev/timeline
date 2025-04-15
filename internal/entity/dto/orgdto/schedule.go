package orgdto

type WorkerSchedule struct {
	OrgID           int         `json:"org_id"`
	WorkerID        int         `json:"worker_id"`
	SessionDuration int         `json:"session_duration,omitempty"`
	Schedule        []*Schedule `json:"schedule"`
}

type ScheduleList struct {
	Workers []*WorkerSchedule `json:"workers"`
	Found   int               `json:"found"`
}

type ScheduleParams struct {
	WorkerID int
	OrgID    int
	Weekday  int
	Limit    int
	Page     int
}

type Schedule struct {
	WorkerScheduleID int    `json:"worker_schedule_id,omitempty"`
	Weekday          int    `json:"weekday"`
	Start            string `json:"start" validate:"time"`
	Over             string `json:"over" validate:"time"`
}
