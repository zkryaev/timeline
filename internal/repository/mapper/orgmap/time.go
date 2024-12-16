package orgmap

import (
	"database/sql"
	"time"
	"timeline/internal/entity"
	"timeline/internal/repository/models/orgmodel"
)

const (
	dateFormat = "01-01-2001"
	timeFormat = "15:04"
)

func OpenHoursToModel(day *entity.OpenHours) *orgmodel.OpenHours {
	open, _ := time.Parse(timeFormat, day.Open)
	close, _ := time.Parse(timeFormat, day.Close)
	breakstart, _ := time.Parse(timeFormat, day.BreakStart)
	breakend, _ := time.Parse(timeFormat, day.BreakEnd)

	return &orgmodel.OpenHours{
		Weekday:    sql.NullInt32{Int32: int32(day.Weekday), Valid: true},
		Open:       sql.NullTime{Time: open, Valid: true},
		Close:      sql.NullTime{Time: close, Valid: true},
		BreakStart: sql.NullTime{Time: breakstart, Valid: true},
		BreakEnd:   sql.NullTime{Time: breakend, Valid: true},
	}
}

func OpenHoursToDTO(day *orgmodel.OpenHours) *entity.OpenHours {
	// Если weekday = -1 значит пустая структура
	if !day.Weekday.Valid {
		return nil
	}
	return &entity.OpenHours{
		Weekday:    int(day.Weekday.Int32),
		Open:       day.Open.Time.Format(timeFormat),
		Close:      day.Close.Time.Format(timeFormat),
		BreakStart: day.BreakStart.Time.Format(timeFormat),
		BreakEnd:   day.BreakEnd.Time.Format(timeFormat),
	}
}
