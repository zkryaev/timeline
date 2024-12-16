package custom

import "time"

// Сравнение по времени (формат 15:00):
// если t < u, то -1,
// если t > u, то +1,
// если t = u, то 0,
func CompareTime(t time.Time, u time.Time) int {
	switch {
	// t = u
	case t.Hour() == u.Hour(), t.Minute() == u.Minute():
		return 0
	// t > u
	case (t.Hour() < u.Hour() && t.Minute() < u.Minute()) || (t.Hour() == u.Hour() && t.Minute() < u.Minute()):
		return -1
	// t < u
	case (t.Hour() > u.Hour() && t.Minute() > u.Minute()) || (t.Hour() == u.Hour() && t.Minute() > u.Minute()):
		return 1
	}
	return 0
}
