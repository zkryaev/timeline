package common

import (
	"timeline/internal/usecase/auth"
	"timeline/internal/usecase/common"
)

// common
var (
	ErrNotFound       = common.ErrNotFound
	ErrNothingChanged = common.ErrNothingChanged
)

// auth
var (
	ErrFailedLogin     = "invalid username or password"
	ErrFailedRegister  = "register failed"
	ErrAccountNotFound = auth.ErrAccountNotFound
	ErrAccountExpired  = auth.ErrAccountExpired
	ErrCodeExpired     = auth.ErrCodeExpired
)

var (
	ErrTimeIncorrect = common.ErrTimeIncorrect
)
