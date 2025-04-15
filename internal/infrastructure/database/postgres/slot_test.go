package postgres

import (
	"context"
	"fmt"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"
)

func (suite *PostgresTestSuite) TestSlotQueries() {
	ctx := context.Background()
	params := &orgdto.SlotReq{
		WorkerID: 1,
		OrgID:    1,
		TData:    entity.TokenData{ID: 1, IsOrg: true},
	}
	slots, _, err := suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
	suite.Require().NoError(err, fmt.Sprintf("worker_id=%d org_id=%d", params.WorkerID, params.OrgID))
	suite.Require().NotNil(slots)

	numSlotsBefore := len(slots)
	suite.NoError(suite.db.GenerateSlots(ctx), "generate slots failed")
	rawNumSlots, _, err := suite.db.Slots(ctx, orgmap.SlotReqToModel(params))
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
