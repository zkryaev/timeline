package models

import "time"

type ExpInfo struct {
	ID        int
	Verified  bool
	CreatedAt time.Time
	Hash      string
}

type CodeInfo struct {
	ID    int
	Code  string
	IsOrg bool
}

type HashCreds struct {
	Email      string `db:"email"`
	PasswdHash string `db:"passwd_hash"`
}