package postgres

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/mapper/recordmap"
)

func (suite *PostgresTestSuite) TestRecordQueries() {
	ctx := context.Background()

	org, err := suite.db.OrgByID(ctx, 1)
	suite.Require().NoError(err)
	suite.Require().NotNil(org)

	user, err := suite.db.UserByID(ctx, 1)
	suite.Require().NoError(err)
	suite.Require().NotNil(user)

	dbWorkers, found, err := suite.db.WorkerList(ctx, org.OrgID, 5, 0)
	suite.Require().NoError(err)
	suite.Require().Greater(found, 0)
	suite.Require().NotNil(dbWorkers)

	params := &orgdto.SlotReq{
		WorkerID: dbWorkers[0].WorkerID,
		OrgID:    org.OrgID,
	}
	dbSlots, err := suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
	suite.Require().NoError(err)
	suite.Require().Greater(len(dbSlots), 0)
	suite.Require().NotNil(dbSlots)

	var freeSlot int
	var isFound bool
	for i := range dbSlots {
		if !dbSlots[i].Busy && dbSlots[i].Date.Day() > time.Now().Day() {
			freeSlot = i
			isFound = true
			break
		}
	}
	suite.Require().True(isFound, "needed free slot not found")

	dbServices, found, err := suite.db.ServiceList(ctx, org.OrgID, 5, 0)
	suite.Require().NoError(err)
	suite.Require().Greater(found, 0)
	suite.Require().NotNil(dbServices)

	addReq := &recordto.Record{
		OrgID:     org.OrgID,
		SlotID:    dbSlots[freeSlot].SlotID,
		UserID:    user.UserID,
		WorkerID:  dbWorkers[0].WorkerID,
		ServiceID: dbServices[0].ServiceID,
	}
	remindRec, recordID, err := suite.db.RecordAdd(ctx, recordmap.RecordToModel(addReq))
	suite.NoError(err, fmt.Sprintf("org_id=%d slot_id=%d user_id=%d worker_id=%d service_id=%d",
		org.OrgID, dbSlots[freeSlot].SlotID, user.UserID, dbWorkers[0].WorkerID, dbServices[0].ServiceID))
	suite.Greater(recordID, 0)
	suite.NotNil(remindRec)

	record, err := suite.db.Record(ctx, recordID)
	suite.NoError(err, fmt.Sprintf("record_id=%d", recordID))
	suite.NotNil(record)

	userAnother, err := suite.db.UserByID(ctx, 2)
	suite.Require().NoError(err)
	suite.Require().NotNil(user)

	patchReq := &recordto.Record{
		RecordID: recordID,
		UserID:   userAnother.UserID,
	}
	suite.NoError(suite.db.RecordPatch(ctx, recordmap.RecordToModel(patchReq)), fmt.Sprintf("record_id=%d user_id=%d", patchReq.RecordID, patchReq.UserID))

	recParams := &recordto.RecordListParams{
		OrgID: org.OrgID,
		Fresh: true,
		Limit: 10,
		Page:  1,
	}
	dbRecords, found, err := suite.db.RecordList(ctx, recordmap.RecordParamsToModel(recParams))
	suite.NoError(err, fmt.Sprintf("org_id=%d fresh=%t limit=%d page=%d", recParams.OrgID, recParams.Fresh, recParams.Limit, recParams.Page))
	suite.Greater(found, 0)
	suite.NotNil(dbRecords)
	for i := range dbRecords {
		if dbRecords[i].RecordID == patchReq.RecordID {
			suite.Equal(userAnother.FirstName, dbRecords[i].User.FirstName)
			suite.Equal(userAnother.LastName, dbRecords[i].User.LastName)
			break
		}
	}
	cancelReq := &recordto.RecordCancelation{
		RecordID: recordID,
		CancelReason: "TESTING REASON",
	}
	suite.NoError(suite.db.RecordCancel(ctx, recordmap.CancelationToModel(cancelReq)), fmt.Sprintf("record_id=%d | slot_id=%d slot.date=%s slot.begin=%s slot-end=%s",
		recordID, dbSlots[freeSlot].SlotID, dbSlots[freeSlot].Date.String(), dbSlots[freeSlot].Begin.Format("15:04"), dbSlots[freeSlot].End.Format("15:04")))
	suite.db.RecordDelete(ctx, recordID)
}
