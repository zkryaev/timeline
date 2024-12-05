package orgdto

type ScheduleList struct {
	WorkerID        int         `json:"worker_id"`
	OrgID           int         `json:"org_id"`
	SessionDuration int         `json:"session_duration,omitempty"`
	Schedule        []*Schedule `json:"schedule"`
}

type ScheduleParams struct {
	WorkerID int `json:"worker_id"`
	OrgID    int `json:"org_id"`
	Weekday  int `json:"weekday"`
}

type Schedule struct {
	WorkerScheduleID int    `json:"worker_schedule_id,omitempty"`
	Weekday          int    `json:"weekday"`
	Start            string `json:"start" validate:"time"`
	Over             string `json:"over" validate:"time"`
}
