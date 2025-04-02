package common

import (
	"errors"
)

var (
	ErrNothingChanged = errors.New("nothing was done")
	ErrNotFound       = errors.New("resource not found")
)
