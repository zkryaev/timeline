package validation

import "time"

// Если over >= start - false
func IsPeriodValid(start, over string) bool {
	overTime, errover := time.Parse("15:06", over)
	startTime, errstart := time.Parse("15:06", start)
	if overTime.Compare(startTime) <= 0 || errover != nil || errstart != nil {
		return false
	}
	return true
}
