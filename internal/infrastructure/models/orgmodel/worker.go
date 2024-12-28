package orgmodel

type Worker struct {
	WorkerID        int    `db:"worker_id"`
	UUID            string `db:"uuid"`
	OrgID           int    `db:"org_id"`
	FirstName       string `db:"first_name"`
	LastName        string `db:"last_name"`
	Position        string `db:"position"`
	Degree          string `db:"degree"`
	SessionDuration int    `db:"session_duration"`
}

type WorkerAssign struct {
	WorkerID  int `db:"worker_id"`
	OrgID     int `db:"org_id"`
	ServiceID int `db:"service_id"`
}
