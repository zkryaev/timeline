package recordmodel

import (
	"database/sql"
	"timeline/internal/infrastructure/models"
)

type Feedback struct {
	FeedbackID      sql.NullInt32  `db:"feedback_id"`
	RecordID        sql.NullInt32  `db:"record_id"`
	Stars           sql.NullInt32  `db:"stars"`
	Feedback        sql.NullString `db:"feedback"`
	Service         sql.NullString `db:"service_name"`
	WorkerFirstName sql.NullString `db:"worker_first_name"`
	WorkerLastName  sql.NullString `db:"worker_last_name"`
	UserFirstName   sql.NullString `db:"user_first_name"`
	UserLastName    sql.NullString `db:"user_last_name"`
	RecordDate      sql.NullTime   `db:"record_date"`
	TData           models.TokenData
}

type FeedbackParams struct {
	TData      models.TokenData
	FeedbackID int `db:"feedback_id"`
	RecordID   int `db:"record_id"`
	UserID     int `db:"user_id"`
	OrgID      int `db:"org_id"`
	Limit      int
	Offset     int
}
