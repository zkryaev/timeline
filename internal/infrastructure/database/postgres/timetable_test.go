package postgres

import (
	"context"
	"time"
	"timeline/internal/entity"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/orgmodel"
)

func IsTimesEqual(suite *PostgresTestSuite, exp, actual *entity.OpenHours) {

}

func (suite *PostgresTestSuite) TestTimetableQueries() {
	ctx := context.Background()

	// (1, 1, '2024-11-28 08:00:00', '2024-11-28 20:00:00', '2024-11-28 12:00:00', '2024-11-28 13:00:00')

	orgID := 1
	modelTData := models.TokenData{ID: orgID, IsOrg: true}
	exp := &entity.OpenHours{
		Weekday:    1,
		Open:       "08:00",
		Close:      "20:00",
		BreakStart: "12:00",
		BreakEnd:   "13:00",
	}
	suite.NoError(suite.db.TimetableAdd(ctx, orgID, []*orgmodel.OpenHours{orgmap.OpenHoursToModel(exp)}))
	timetable, _, err := suite.db.Timetable(ctx, orgID, modelTData)
	suite.NoError(err)
	suite.NotNil(timetable)

	openhours := orgmap.TimetableToEntity(timetable, time.Now().Location())
	var found bool
	for _, t := range openhours {
		if t.Weekday == exp.Weekday {
			suite.Equal(exp.Weekday, t.Weekday)
			found = true
		}
	}
	suite.NotZero(found)

	exp.Open = "10:00"
	exp.Close = "19:00"
	suite.NoError(suite.db.TimetableUpdate(ctx, orgID, []*orgmodel.OpenHours{orgmap.OpenHoursToModel(exp)}))

	timetable, _, err = suite.db.Timetable(ctx, orgID, modelTData)
	suite.NoError(err)
	suite.NotNil(timetable)

	found = false
	openhours = orgmap.TimetableToEntity(timetable, time.Now().UTC().Location())
	for _, t := range openhours {
		if t.Weekday == exp.Weekday {
			suite.Equal(exp, t)
			found = true
		}
	}
	suite.NotZero(found)

	suite.NoError(suite.db.TimetableDelete(ctx, orgID, exp.Weekday))

	timetable, _, err = suite.db.Timetable(ctx, orgID, modelTData)
	suite.NoError(err)
	suite.NotNil(timetable)

	found = false
	openhours = orgmap.TimetableToEntity(timetable, time.Now().Location())
	for _, t := range openhours {
		if t.Weekday == exp.Weekday {
			found = true
		}
	}
	suite.Zero(found)
}
