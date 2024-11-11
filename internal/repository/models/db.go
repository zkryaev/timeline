package models

import "time"

type IsExistResponse struct {
	ID        int
	Hash      string
	CreatedAt time.Time
}
