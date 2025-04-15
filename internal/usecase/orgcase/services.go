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

func (o *OrgUseCase) Service(ctx context.Context, logger *zap.Logger, serviceID, orgID int) (*orgdto.ServiceList, error) {
	service, err := o.org.Service(ctx, serviceID, orgID)
	if err != nil {
		if errors.Is(err, postgres.ErrServiceNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched service")
	list := &orgdto.ServiceList{List: []*orgdto.ServiceResp{orgmap.ServiceToDTO(service)}, Found: 1}
	return list, nil
}

func (o *OrgUseCase) ServiceList(ctx context.Context, logger *zap.Logger, orgID int, limit int, page int) (*orgdto.ServiceList, error) {
	offset := (page - 1) * limit
	data, found, err := o.org.ServiceList(ctx, orgID, limit, offset)
	if err != nil {
		if errors.Is(err, postgres.ErrServiceNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched service list")
	serviceList := make([]*orgdto.ServiceResp, 0, len(data))
	for _, v := range data {
		serviceList = append(serviceList, orgmap.ServiceToDTO(v))
	}
	resp := &orgdto.ServiceList{
		List:  serviceList,
		Found: found,
	}
	return resp, nil
}

func (o *OrgUseCase) ServiceAdd(ctx context.Context, logger *zap.Logger, service *orgdto.AddServiceReq) error {
	_, err := o.org.ServiceAdd(ctx, orgmap.AddServiceToModel(service))
	if err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Service has been saved")
	return nil
}

func (o *OrgUseCase) ServiceUpdate(ctx context.Context, logger *zap.Logger, service *orgdto.UpdateServiceReq) error {
	if err := o.org.ServiceUpdate(ctx, orgmap.UpdateService(service)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Service has been updated")
	return nil
}

func (o *OrgUseCase) ServiceDelete(ctx context.Context, logger *zap.Logger, serviceID, orgID int) error {
	if err := o.org.ServiceSoftDelete(ctx, serviceID, orgID); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Service has been deleted")
	return nil
}
