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

func (o *OrgUseCase) Worker(ctx context.Context, logger *zap.Logger, workerID, orgID int) (*orgdto.WorkerList, error) {
	worker, err := o.org.Worker(ctx, workerID, orgID)
	if err != nil {
		if errors.Is(err, postgres.ErrWorkerNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched worker")
	workerList := []*orgdto.WorkerResp{orgmap.WorkerToDTO(worker)}
	resp := &orgdto.WorkerList{
		List:  workerList,
		Found: 1,
	}
	return resp, nil
}

func (o *OrgUseCase) WorkerList(ctx context.Context, logger *zap.Logger, orgID, limit, page int) (*orgdto.WorkerList, error) {
	offset := (page - 1) * limit
	data, found, err := o.org.WorkerList(ctx, orgID, limit, offset)
	if err != nil {
		if errors.Is(err, postgres.ErrWorkerNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched worker list")
	workerList := make([]*orgdto.WorkerResp, 0, len(data))
	for _, v := range data {
		workerList = append(workerList, orgmap.WorkerToDTO(v))
	}
	resp := &orgdto.WorkerList{
		List:  workerList,
		Found: found,
	}
	return resp, nil
}

func (o *OrgUseCase) WorkerAdd(ctx context.Context, logger *zap.Logger, worker *orgdto.AddWorkerReq) (*orgdto.WorkerResp, error) {
	workerID, err := o.org.WorkerAdd(ctx, orgmap.AddWorkerToModel(worker))
	if err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return nil, common.ErrNothingChanged
		}
		return nil, err
	}
	logger.Info("Worker has been saved")
	return &orgdto.WorkerResp{
		WorkerID: workerID,
	}, nil
}

func (o *OrgUseCase) WorkerUpdate(ctx context.Context, logger *zap.Logger, worker *orgdto.UpdateWorkerReq) error {
	if err := o.org.WorkerUpdate(ctx, orgmap.UpdateWorkerToModel(worker)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Worker has been updated")
	return nil
}

func (o *OrgUseCase) WorkerDelete(ctx context.Context, logger *zap.Logger, workerID, orgID int) error {
	if err := o.org.WorkerSoftDelete(ctx, workerID, orgID); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Worker has been deleted")
	return nil
}

func (o *OrgUseCase) WorkersServices(ctx context.Context, logger *zap.Logger, serviceID, orgID int) ([]*orgdto.WorkerResp, error) {
	data, err := o.org.ServiceWorkerList(ctx, serviceID, orgID)
	if err != nil {
		if errors.Is(err, postgres.ErrServiceNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched worker-service list")
	workers := make([]*orgdto.WorkerResp, 0, len(data))
	for _, worker := range data {
		workers = append(workers, orgmap.WorkerToDTO(worker))
	}
	return workers, nil
}

func (o *OrgUseCase) WorkerAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error {
	if err := o.org.WorkerAssignService(ctx, orgmap.AssignWorkerToModel(assignInfo)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Worker has been assigned to service")
	return nil
}

func (o *OrgUseCase) WorkerUnAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error {
	if err := o.org.WorkerUnAssignService(ctx, orgmap.AssignWorkerToModel(assignInfo)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Worker has been unassigned from service")
	return nil
}
