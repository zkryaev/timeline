package postgres

import (
	"context"
	"fmt"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"
)

func (suite *PostgresTestSuite) TestScheduleQueries() {
	ctx := context.Background()

	org, err := suite.db.OrgByID(ctx, 1)
	suite.Require().NoError(err)
	suite.Require().NotNil(org)

	dbWorkers, found, err := suite.db.WorkerList(ctx, org.OrgID, 5, 0)
	suite.Require().NoError(err)
	suite.Require().Greater(found, 0)
	suite.Require().NotNil(dbWorkers)

	worker := orgmap.WorkerToDTO(dbWorkers[0])

	params := &orgdto.ScheduleParams{
		WorkerID: worker.WorkerID,
		OrgID:    org.OrgID,
		Weekday:  1,
		Limit:    10,
		Page:     1,
	}
	schedules, err := suite.db.WorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params))
	suite.NoError(err, fmt.Sprintf("worker_id=%d org_id=%d weekday=%d limit=%d page=%d",
		params.WorkerID, params.OrgID, params.Weekday, params.Limit, params.Page))
	suite.NotNil(schedules)

	addReq := &orgdto.WorkerSchedule{
		WorkerID:        worker.WorkerID,
		OrgID:           org.OrgID,
		SessionDuration: schedules.Workers[0].SessionDuration,
		Schedule: []*orgdto.Schedule{
			{
				Weekday: 3,
				Start:   "14:00",
				Over:    "18:00",
			},
		},
	}
	suite.NoError(suite.db.AddWorkerSchedule(ctx, orgmap.WorkerScheduleToModel(addReq)),
		fmt.Sprintf("weekday=%d start=%s over=%s", addReq.Schedule[0].Weekday, addReq.Schedule[0].Start, addReq.Schedule[0].Over))

	params.Weekday = addReq.Schedule[0].Weekday
	scheduleList, err := suite.db.WorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params))
	suite.NoError(err, fmt.Sprintf("worker_id=%d org_id=%d weekday=%d limit=%d page=%d",
		params.WorkerID, params.OrgID, params.Weekday, params.Limit, params.Page))
	suite.NotNil(scheduleList)
	var isFound bool
	var indExpected [2]int
	for i, workerSched := range scheduleList.Workers {
		if workerSched.WorkerID == worker.WorkerID && workerSched.OrgID == org.OrgID {
			scheds := orgmap.WorkerScheduleToDTO(workerSched)
			for j, sched := range scheds.Schedule {
				if addReq.Schedule[0].Weekday == sched.Weekday && addReq.Schedule[0].Start == sched.Start && addReq.Schedule[0].Over == sched.Over {
					isFound = true
					indExpected[0] = i
					indExpected[1] = j
					break
				}
			}
		}
	}
	suite.Require().True(isFound, "added schedule not found")

	addReq.Schedule[0].WorkerScheduleID = scheduleList.Workers[indExpected[0]].Schedule[indExpected[1]].WorkerScheduleID
	addReq.Schedule[0].Weekday = 4
	suite.NoError(suite.db.UpdateWorkerSchedule(ctx, orgmap.WorkerScheduleToModel(addReq)),
		fmt.Sprintf("weekday=%d start=%s over=%s", addReq.Schedule[0].Weekday, addReq.Schedule[0].Start, addReq.Schedule[0].Over))

	scheduleList, err = suite.db.WorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params))
	suite.NoError(err, fmt.Sprintf("worker_id=%d org_id=%d weekday=%d limit=%d page=%d",
		params.WorkerID, params.OrgID, params.Weekday, params.Limit, params.Page))
	suite.NotNil(scheduleList)
	for _, workerSched := range scheduleList.Workers {
		if workerSched.WorkerID == worker.WorkerID && workerSched.OrgID == org.OrgID {
			scheds := orgmap.WorkerScheduleToDTO(workerSched)
			for _, sched := range scheds.Schedule {
				if addReq.Schedule[0].Weekday == sched.Weekday && addReq.Schedule[0].Start == sched.Start && addReq.Schedule[0].Over == sched.Over {
					break
				}
			}
		}
	}
	schedparams := &orgdto.ScheduleParams{
		WorkerID: worker.WorkerID,
		OrgID:    org.OrgID,
		Weekday:  addReq.Schedule[0].Weekday,
	}
	suite.NoError(suite.db.DeleteWorkerSchedule(ctx, orgmap.ScheduleParamsToModel(schedparams)), "couldn't delete created schedule")
}
