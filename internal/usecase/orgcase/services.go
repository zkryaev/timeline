package orgcase

import (
	"context"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) Service(ctx context.Context, serviceID, orgID int) (*orgdto.ServiceResp, error) {
	service, err := o.org.Service(ctx, serviceID, orgID)
	if err != nil {
		o.Logger.Error(
			"failed to get service",
			zap.Error(err),
		)
		return nil, err
	}
	return orgmap.ServiceToDTO(service), nil
}

func (o *OrgUseCase) ServiceWorkerList(ctx context.Context, serviceID, orgID int) ([]*orgdto.WorkerResp, error) {
	data, err := o.org.ServiceWorkerList(ctx, serviceID, orgID)
	if err != nil {
		o.Logger.Error(
			"failed to get worker service list",
			zap.Error(err),
		)
		return nil, err
	}
	workers := make([]*orgdto.WorkerResp, 0, len(data))
	for _, worker := range data {
		workers = append(workers, orgmap.WorkerToDTO(worker))
	}

	return workers, nil
}
func (o *OrgUseCase) ServiceAdd(ctx context.Context, service *orgdto.AddServiceReq) error {
	_, err := o.org.ServiceAdd(ctx, orgmap.AddServiceToModel(service))
	if err != nil {
		o.Logger.Error(
			"failed to add service",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) ServiceUpdate(ctx context.Context, service *orgdto.UpdateServiceReq) error {
	if err := o.org.ServiceUpdate(ctx, orgmap.UpdateService(service)); err != nil {
		o.Logger.Error(
			"failed to update service",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) ServiceList(ctx context.Context, orgID int, limit int, page int) (*orgdto.ServiceList, error) {
	offset := (page - 1) * limit
	data, found, err := o.org.ServiceList(ctx, orgID, limit, offset)
	if err != nil {
		o.Logger.Error(
			"failed to retrieve list of services",
			zap.Error(err),
		)
		return nil, nil //nolint:nilnil // correct
	}
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

func (o *OrgUseCase) ServiceDelete(ctx context.Context, serviceID, orgID int) error {
	if err := o.org.ServiceSoftDelete(ctx, serviceID, orgID); err != nil {
		o.Logger.Error(
			"failed to delete service",
			zap.Error(err),
		)
	}
	return nil
}
