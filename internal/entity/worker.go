package entity

type Worker struct {
	UUID            string   `json:"uuid"`
	FirstName       string   `json:"first_name,omitempty"`
	LastName        string   `json:"last_name,omitempty"`
	Position        string   `json:"position,omitempty"`
	Degree          string   `json:"degree,omitempty"`
	SessionDuration int      `json:"session_duration,omitempty"`
	ImagesURL       []string `json:"images_url,omitempty"`
}
