package models

import (
	"context"
	"time"
)

type ReminderMsg struct {
	Organization string
	Service      string
	ServiceDesc  string
	Address      string
	SessionStart time.Time
	SessionEnd   time.Time
	SessionDate  time.Time
}

type CancelMsg struct {
	Organization string
	Service      string
	ServiceDecs  string
	SessionStart string
	SessionEnd   string
	SessionDate  string
	CancelReason string
}

type Message struct {
	Email    string
	Type     string
	Value    interface{}
	IsAttach bool
}

type WorkerContext struct {
	Context context.Context
	Close   context.CancelFunc
}
