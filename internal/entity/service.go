package entity

type Service struct {
	Name        string  `json:"name,omitempty"`
	Cost        float64 `json:"cost,omitempty"`
	Description string  `json:"description,omitempty"`
}
