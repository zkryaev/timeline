package custom

import "time"

// Сравнение по времени (формат 15:00):
// если t < u, то -1,
// если t > u, то +1,
// если t = u, то 0,
func CompareTime(t time.Time, u time.Time) int {
	tHour, tMinute := t.Hour(), t.Minute()
	uHour, uMinute := u.Hour(), u.Minute()

	if tHour < uHour {
		return -1
	}
	if tHour > uHour {
		return 1
	}
	if tMinute < uMinute {
		return -1
	}
	if tMinute > uMinute {
		return 1
	}
	return 0
}
