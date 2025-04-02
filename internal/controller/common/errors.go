package common

import "timeline/internal/usecase/common"

// auth.go
var (
	ErrFailedLogin    = "invalid username or password"
	ErrFailedRegister = "register failed"
	ErrNotFound       = common.ErrNotFound
	ErrNothingChanged = common.ErrNothingChanged
)
