package models

import (
	"context"
	"time"
)

type ReminderMsg struct {
	Organization string
	Service      string
	Description  string
	Address      string
	SessionStart time.Time
	SessionEnd   time.Time
	SessionDate  time.Time
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
