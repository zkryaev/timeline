package recordto

type Feedback struct {
	FeedbackID int    `json:"feedback_id,omitempty"`
	RecordID   int    `json:"record_id,omitempty"`
	Stars      int    `json:"stars,omitempty"`
	Feedback   string `json:"feedback,omitempty"`
}

type FeedbackParams struct {
	FeedbackID int `json:"feedback_id"`
	RecordID   int `json:"record_id"`
}
