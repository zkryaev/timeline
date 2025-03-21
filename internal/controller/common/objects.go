package common

import (
	"timeline/internal/controller/validation"

	jsoniter "github.com/json-iterator/go"
)

// singleton object
var (
	fastjson     = jsoniter.ConfigFastest // float with 6 digits, no sort
	validator, _ = validation.NewCustomValidator()
)
