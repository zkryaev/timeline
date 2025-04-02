package orgcase

import (
	"context"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) UpdateSlot(ctx context.Context, logger *zap.Logger, req *orgdto.SlotUpdate) error {
	if err := o.org.UpdateSlot(ctx, req.Busy, orgmap.SlotUpdateToModel(req)); err != nil {
		return err
	}
	logger.Info("Slot has been updated")
	return nil
}

func (o *OrgUseCase) Slots(ctx context.Context, logger *zap.Logger, req *orgdto.SlotReq) ([]*orgdto.SlotResp, error) {
	data, err := o.org.Slots(ctx, orgmap.SlotReqToModel(req))
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched slots")
	resp := make([]*orgdto.SlotResp, 0, len(data))
	for _, v := range data {
		resp = append(resp, orgmap.SlotToDTO(v))
	}
	return resp, nil
}
