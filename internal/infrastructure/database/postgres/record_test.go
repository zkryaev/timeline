//go:build integration

package postgres

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/mapper/recordmap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/recordmodel"
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
		TData:    entity.TokenData{ID: org.OrgID, IsOrg: true},
	}
	dbSlots, _, err := suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
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
	suite.Require().NotNil(remindRec)

	req := recordmodel.RecordParam{RecordID: recordID, TData: models.TokenData{ID: user.UserID, IsOrg: false}}
	record, err := suite.db.Record(ctx, req)
	suite.NoError(err, fmt.Sprintf("record_id=%d", recordID))
	suite.Require().NotNil(record)

	recParams := &recordto.RecordListParams{}
	recParams.OrgID = 0
	recParams.UserID = 1
	recParams.Fresh = false
	recParams.Page = 1
	_, found, err = suite.db.RecordList(ctx, recordmap.RecordParamsToModel(recParams))
	suite.NoError(err)
	suite.Greater(found, 0)

	cancelReq := &recordto.RecordCancelation{
		TData:        entity.TokenData{ID: user.UserID},
		RecordID:     recordID,
		CancelReason: "TESTING REASON",
	}
	suite.NoError(suite.db.RecordCancel(ctx, recordmap.CancelationToModel(cancelReq)), fmt.Sprintf("record_id=%d | slot_id=%d slot.date=%s slot.begin=%s slot-end=%s",
		recordID, dbSlots[freeSlot].SlotID, dbSlots[freeSlot].Date.String(), dbSlots[freeSlot].Begin.Format("15:04"), dbSlots[freeSlot].End.Format("15:04")))
	suite.db.RecordDelete(ctx, recordID)
}
