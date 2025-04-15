package orgcase

import (
	"context"
	"errors"
	"time"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/usecase/common"

	"go.uber.org/zap"
)

func (o *OrgUseCase) Slots(ctx context.Context, logger *zap.Logger, req *orgdto.SlotReq) ([]*orgdto.SlotResp, error) {
	data, city, err := o.org.Slots(ctx, orgmap.SlotReqToModel(req))
	if err != nil {
		if errors.Is(err, postgres.ErrSlotsNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	tzid := o.backdata.Cities.GetCityTZ(city)
	loc, err := time.LoadLocation(tzid)
	if err != nil {
		logger.Error("failed to load location, set UTC+03 (MSK)", zap.String("city-tzid", city+"="+tzid), zap.Error(err))
		loc = time.Local // UTC+03 = MSK
	}
	logger.Info("Fetched slots")
	resp := make([]*orgdto.SlotResp, 0, len(data))
	for _, v := range data {
		resp = append(resp, orgmap.SlotToDTO(v, loc))
	}
	return resp, nil
}
