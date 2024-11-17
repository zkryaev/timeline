package models

import "time"

type MetaInfo struct {
	ID        int
	Hash      string
	CreatedAt time.Time
	Verified  bool
}

type HashCreds struct {
	Email      string `db:"email"`
	PasswdHash string `db:"passwd_hash"`
}
