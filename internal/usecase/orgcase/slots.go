package orgcase

import (
	"context"
	"errors"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/usecase/common"

	"go.uber.org/zap"
)

func (o *OrgUseCase) UpdateSlot(ctx context.Context, logger *zap.Logger, req *orgdto.SlotUpdate) error {
	if err := o.org.UpdateSlot(ctx, req.Busy, orgmap.SlotUpdateToModel(req)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Slot has been updated")
	return nil
}

func (o *OrgUseCase) Slots(ctx context.Context, logger *zap.Logger, req *orgdto.SlotReq) ([]*orgdto.SlotResp, error) {
	data, err := o.org.Slots(ctx, orgmap.SlotReqToModel(req))
	if err != nil {
		if errors.Is(err, postgres.ErrSlotsNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched slots")
	resp := make([]*orgdto.SlotResp, 0, len(data))
	for _, v := range data {
		resp = append(resp, orgmap.SlotToDTO(v))
	}
	return resp, nil
}
