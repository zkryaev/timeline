package custom

import "time"

// Сравнение только по часам и минутам (формат 15:00):
//   - если a < b, то -1
//   - если a > b, то +1
//   - если a = b, то  0
func CompareTime(a time.Time, b time.Time) int {
	aHour, aMinute := a.Hour(), a.Minute()
	bHour, bMinute := b.Hour(), b.Minute()

	if aHour < bHour {
		return -1
	}
	if aHour > bHour {
		return 1
	}
	if aMinute < bMinute {
		return -1
	}
	if aMinute > bMinute {
		return 1
	}
	return 0
}
