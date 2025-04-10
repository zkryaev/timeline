package postgres

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"
)

func (suite *PostgresTestSuite) TestSlotQueries() {
	ctx := context.Background()
	params := &orgdto.SlotReq{
		WorkerID: 1,
		OrgID:    1,
	}
	slots, _, err := suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
	suite.Require().NoError(err, fmt.Sprintf("worker_id=%d org_id=%d", params.WorkerID, params.OrgID))
	suite.NotNil(slots)

	freeSlot := &orgdto.SlotResp{}

	for _, slot := range slots {
		if !slot.Busy {
			freeSlot.Slot = *orgmap.SlotInfoToDTO(slot, time.Now().Location())
			freeSlot.SlotID = slot.SlotID
			break
		}
	}
	updparams := &orgdto.SlotUpdate{
		SlotID:   freeSlot.SlotID,
		WorkerID: freeSlot.WorkerID,
	}
	suite.NoError(suite.db.UpdateSlot(ctx, true, orgmap.SlotUpdateToModel(updparams)), fmt.Sprintf("slot_id=%d worker_id=%d", freeSlot.SlotID, freeSlot.WorkerID))

	params.UserID = 1
	params.OrgID = 0
	params.WorkerID = 0
	rawNumSlots, _, err := suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
	suite.Require().NoError(err)
	suite.Greater(len(rawNumSlots), 0)
	numSlotsBefore := len(rawNumSlots)
	suite.NoError(suite.db.GenerateSlots(ctx), "generate slots failed")
	rawNumSlots, _, err = suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
	suite.Require().NoError(err)
	suite.Greater(len(rawNumSlots), 0)
	numSlotsAfterGenerations := len(rawNumSlots)
	suite.Greater(numSlotsAfterGenerations, numSlotsBefore)

	suite.NoError(suite.db.DeleteExpiredSlots(ctx))
	rawNumSlots, _, err = suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
	suite.Require().NoError(err)
	suite.Greater(len(rawNumSlots), 0)
	numSlotsAfterDeletions := len(rawNumSlots)
	suite.Less(numSlotsAfterDeletions, numSlotsAfterGenerations)
}
