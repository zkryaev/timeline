package recordmodel

import "database/sql"

type Feedback struct {
	FeedbackID sql.NullInt32  `db:"feedback_id"`
	RecordID   sql.NullInt32  `db:"record_id"`
	Stars      sql.NullInt32  `db:"stars"`
	Feedback   sql.NullString `db:"feedback"`
}

type FeedbackParams struct {
	FeedbackID int `db:"feedback_id"`
	RecordID   int `db:"record_id"`
}
