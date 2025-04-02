package orgcase

import (
	"context"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) Worker(ctx context.Context, logger *zap.Logger, workerID, orgID int) (*orgdto.WorkerResp, error) {
	worker, err := o.org.Worker(ctx, workerID, orgID)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched worker")
	return orgmap.WorkerToDTO(worker), nil
}
func (o *OrgUseCase) WorkerAdd(ctx context.Context, logger *zap.Logger, worker *orgdto.AddWorkerReq) (*orgdto.WorkerResp, error) {
	workerID, err := o.org.WorkerAdd(ctx, orgmap.AddWorkerToModel(worker))
	if err != nil {
		return nil, err
	}
	logger.Info("Worker has been saved")
	return &orgdto.WorkerResp{
		WorkerID: workerID,
	}, nil
}
func (o *OrgUseCase) WorkerUpdate(ctx context.Context, logger *zap.Logger, worker *orgdto.UpdateWorkerReq) error {
	if err := o.org.WorkerUpdate(ctx, orgmap.UpdateWorkerToModel(worker)); err != nil {
		return err
	}
	logger.Info("Worker has been updated")
	return nil
}

func (o *OrgUseCase) WorkerPatch(ctx context.Context, logger *zap.Logger, worker *orgdto.UpdateWorkerReq) error {
	if err := o.org.WorkerPatch(ctx, orgmap.UpdateWorkerToModel(worker)); err != nil {
		return err
	}
	logger.Info("Worker has been patched")
	return nil
}

func (o *OrgUseCase) WorkerAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error {
	if err := o.org.WorkerAssignService(ctx, orgmap.AssignWorkerToModel(assignInfo)); err != nil {
		return err
	}
	logger.Info("Worker has been assigned to service")
	return nil
}

func (o *OrgUseCase) WorkerUnAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error {
	if err := o.org.WorkerUnAssignService(ctx, orgmap.AssignWorkerToModel(assignInfo)); err != nil {
		return err
	}
	logger.Info("Worker has been unassigned from service")
	return nil
}

func (o *OrgUseCase) WorkerList(ctx context.Context, logger *zap.Logger, orgID, limit, page int) (*orgdto.WorkerList, error) {
	offset := (page - 1) * limit
	data, found, err := o.org.WorkerList(ctx, orgID, limit, offset)
	if err != nil {
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
func (o *OrgUseCase) WorkerDelete(ctx context.Context, logger *zap.Logger, workerID, orgID int) error {
	if err := o.org.WorkerSoftDelete(ctx, workerID, orgID); err != nil {
		return err
	}
	logger.Info("Worker has been deleted")
	return nil
}
