package orgmap

import (
	"database/sql"
	"time"
	"timeline/internal/entity"
	"timeline/internal/infrastructure/models/orgmodel"
	"timeline/internal/usecase/common"
)

func OpenHoursToModel(day *entity.OpenHours) *orgmodel.OpenHours {
	openTime, _ := time.Parse(common.MinutesOnly, day.Open)
	closeTime, _ := time.Parse(common.MinutesOnly, day.Close)
	breakstart, _ := time.Parse(common.MinutesOnly, day.BreakStart)
	breakend, _ := time.Parse(common.MinutesOnly, day.BreakEnd)

	return &orgmodel.OpenHours{
		Weekday:    sql.NullInt32{Int32: int32(day.Weekday), Valid: true},
		Open:       sql.NullTime{Time: openTime.UTC(), Valid: true},
		Close:      sql.NullTime{Time: closeTime.UTC(), Valid: true},
		BreakStart: sql.NullTime{Time: breakstart.UTC(), Valid: true},
		BreakEnd:   sql.NullTime{Time: breakend.UTC(), Valid: true},
	}
}

func OpenHoursToDTO(day *orgmodel.OpenHours, loc *time.Location) *entity.OpenHours {
	if !day.Weekday.Valid {
		return nil
	}
	return &entity.OpenHours{
		Weekday:    int(day.Weekday.Int32),
		Open:       day.Open.Time.In(loc).Format(common.MinutesOnly),
		Close:      day.Close.Time.In(loc).Format(common.MinutesOnly),
		BreakStart: day.BreakStart.Time.In(loc).Format(common.MinutesOnly),
		BreakEnd:   day.BreakEnd.Time.In(loc).Format(common.MinutesOnly),
	}
}
