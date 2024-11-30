package orgcase

import (
	"context"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) Worker(ctx context.Context, WorkerID, OrgID int) (*orgdto.WorkerResp, error) {
	worker, err := o.org.Worker(ctx, WorkerID, OrgID)
	if err != nil {
		o.Logger.Error(
			"failed to get worker",
			zap.Error(err),
		)
		return nil, err
	}
	return orgmap.WorkerToDTO(worker), nil
}
func (o *OrgUseCase) WorkerAdd(ctx context.Context, worker *orgdto.AddWorkerReq) (*orgdto.WorkerResp, error) {
	workerID, err := o.org.WorkerAdd(ctx, orgmap.AddWorkerToModel(worker))
	if err != nil {
		o.Logger.Error(
			"failed to add worker",
			zap.Error(err),
		)
		return nil, err
	}
	return &orgdto.WorkerResp{
		WorkerID: workerID,
	}, nil
}
func (o *OrgUseCase) WorkerUpdate(ctx context.Context, worker *orgdto.UpdateWorkerReq) (*orgdto.UpdateWorkerReq, error) {
	if err := o.org.WorkerUpdate(ctx, orgmap.UpdateWorkerToModel(worker)); err != nil {
		o.Logger.Error(
			"failed to update worker",
			zap.Error(err),
		)
		return nil, err
	}
	return worker, nil
}
func (o *OrgUseCase) WorkerAssignService(ctx context.Context, assignInfo *orgdto.AssignWorkerReq) error {
	if err := o.org.WorkerAssignService(ctx, orgmap.AssignWorkerToModel(assignInfo)); err != nil {
		o.Logger.Error(
			"failed to assign worker to service",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) WorkerUnAssignService(ctx context.Context, assignInfo *orgdto.AssignWorkerReq) error {
	if err := o.org.WorkerUnAssignService(ctx, orgmap.AssignWorkerToModel(assignInfo)); err != nil {
		o.Logger.Error(
			"failed to assign worker to service",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) WorkerList(ctx context.Context, OrgID int) ([]*orgdto.WorkerResp, error) {
	data, err := o.org.WorkerList(ctx, OrgID)
	if err != nil {
		o.Logger.Error(
			"failed to get worker list",
			zap.Error(err),
		)
		return nil, err
	}
	workerList := make([]*orgdto.WorkerResp, 0, len(data))
	for _, v := range data {
		workerList = append(workerList, orgmap.WorkerToDTO(v))
	}
	return workerList, nil
}
func (o *OrgUseCase) WorkerDelete(ctx context.Context, WorkerID, OrgID int) error {
	if err := o.org.WorkerDelete(ctx, WorkerID, OrgID); err != nil {
		o.Logger.Error(
			"failed to delete worker",
			zap.Error(err),
		)
		return err
	}
	return nil
}
