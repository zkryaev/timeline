package recordmodel

import (
	"time"
	"timeline/internal/infrastructure/models/orgmodel"
	"timeline/internal/infrastructure/models/usermodel"
)

type Record struct {
	RecordID  int  `db:"record_id"`
	OrgID     int  `db:"org_id"`
	UserID    int  `db:"user_id"`
	SlotID    int  `db:"slot_id"`
	ServiceID int  `db:"service_id"`
	WorkerID  int  `db:"worker_id"`
	Reviewed  bool `db:"reviewed"`
}

type RecordListParams struct {
	OrgID    int  `db:"org_id"`
	UserID   int  `db:"user_id"`
	Reviewed bool `db:"reviewed"`
	Fresh    bool
	Limit    int
	Offset   int
}

type RecordScrap struct {
	RecordID  int  `db:"record_id"`
	Reviewed  bool `db:"reviewed"`
	Org       *orgmodel.OrgInfo
	User      *usermodel.UserInfo
	Slot      *orgmodel.Slot
	Service   *orgmodel.Service
	Worker    *orgmodel.Worker
	Feedback  *Feedback
	CreatedAt time.Time
}

type ReminderRecord struct {
	UserEmail          string
	UserCity           string
	ServiceName        string
	ServiceDescription string
	OrgName            string
	OrgAddress         string
	Date               time.Time
	Begin              time.Time
	End                time.Time
}

type RecordCancelation struct {
	RecordID     int
	IsCanceled   bool
	CancelReason string
}
