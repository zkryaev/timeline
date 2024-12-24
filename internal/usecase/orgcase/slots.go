package orgcase

import (
	"context"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) UpdateSlot(ctx context.Context, req *orgdto.SlotUpdate) error {
	if err := o.org.UpdateSlot(ctx, req.Busy, orgmap.SlotUpdateToModel(req)); err != nil {
		o.Logger.Error(
			"failed to update slot",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) Slots(ctx context.Context, req *orgdto.SlotReq) ([]*orgdto.SlotResp, error) {
	data, err := o.org.Slots(ctx, orgmap.SlotReqToModel(req))
	if err != nil {
		o.Logger.Error(
			"failed to get slots",
			zap.Error(err),
		)
		return nil, err
	}
	resp := make([]*orgdto.SlotResp, 0, len(data))
	for _, v := range data {
		resp = append(resp, orgmap.SlotToDTO(v))
	}
	return resp, nil
}
