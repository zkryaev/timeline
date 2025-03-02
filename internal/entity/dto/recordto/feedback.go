package recordto

type Feedback struct {
	FeedbackID      int    `json:"feedback_id"`
	RecordID        int    `json:"record_id"`
	Stars           int    `json:"stars"`
	Feedback        string `json:"feedback,omitempty"`
	Service         string `json:"service_name"`
	WorkerFirstName string `json:"worker_first_name"`
	WorkerLastName  string `json:"worker_last_name"`
	UserFirstName   string `json:"user_first_name"`
	UserLastName    string `json:"user_last_name"`
	RecordDate      string `json:"record_date"`
}

type FeedbackList struct {
	List  []*Feedback `json:"feedback_list"`
	Found int         `json:"found"`
}

type FeedbackParams struct {
	FeedbackID int `json:"feedback_id"`
	RecordID   int `json:"record_id"`
	UserID     int
	OrgID      int
	Limit      int
	Page       int
}
