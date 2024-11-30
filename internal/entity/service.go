package entity

type Service struct {
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Description string  `json:"description"`
}
