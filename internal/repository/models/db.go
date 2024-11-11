package models

import "time"

type MetaInfo struct {
	ID        int
	Hash      string
	CreatedAt time.Time
	Verified  bool
}
