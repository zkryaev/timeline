package entity

type Worker struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Position  string `json:"position,omitempty"`
	Degree    string `json:"Degree,omitempty"`
}
