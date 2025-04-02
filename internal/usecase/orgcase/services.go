package orgcase

import (
	"context"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) Service(ctx context.Context, logger *zap.Logger, serviceID, orgID int) (*orgdto.ServiceResp, error) {
	service, err := o.org.Service(ctx, serviceID, orgID)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched service")
	return orgmap.ServiceToDTO(service), nil
}

func (o *OrgUseCase) ServiceWorkerList(ctx context.Context, logger *zap.Logger, serviceID, orgID int) ([]*orgdto.WorkerResp, error) {
	data, err := o.org.ServiceWorkerList(ctx, serviceID, orgID)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched worker-service list")
	workers := make([]*orgdto.WorkerResp, 0, len(data))
	for _, worker := range data {
		workers = append(workers, orgmap.WorkerToDTO(worker))
	}
	return workers, nil
}
func (o *OrgUseCase) ServiceAdd(ctx context.Context, logger *zap.Logger, service *orgdto.AddServiceReq) error {
	_, err := o.org.ServiceAdd(ctx, orgmap.AddServiceToModel(service))
	if err != nil {
		return err
	}
	logger.Info("Service has been saved")
	return nil
}

func (o *OrgUseCase) ServiceUpdate(ctx context.Context, logger *zap.Logger, service *orgdto.UpdateServiceReq) error {
	if err := o.org.ServiceUpdate(ctx, orgmap.UpdateService(service)); err != nil {
		return err
	}
	logger.Info("Service has been updated")
	return nil
}

func (o *OrgUseCase) ServiceList(ctx context.Context, logger *zap.Logger, orgID int, limit int, page int) (*orgdto.ServiceList, error) {
	offset := (page - 1) * limit
	data, found, err := o.org.ServiceList(ctx, orgID, limit, offset)
	if err != nil {
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

func (o *OrgUseCase) ServiceDelete(ctx context.Context, logger *zap.Logger, serviceID, orgID int) error {
	if err := o.org.ServiceSoftDelete(ctx, serviceID, orgID); err != nil {
		return err
	}
	logger.Info("Service has been deleted")
	return nil
}
