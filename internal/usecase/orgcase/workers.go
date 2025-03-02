package orgcase

import (
	"context"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) Worker(ctx context.Context, workerID, orgID int) (*orgdto.WorkerResp, error) {
	worker, err := o.org.Worker(ctx, workerID, orgID)
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
func (o *OrgUseCase) WorkerUpdate(ctx context.Context, worker *orgdto.UpdateWorkerReq) error {
	if err := o.org.WorkerUpdate(ctx, orgmap.UpdateWorkerToModel(worker)); err != nil {
		o.Logger.Error(
			"failed to update worker",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) WorkerPatch(ctx context.Context, worker *orgdto.UpdateWorkerReq) error {
	if err := o.org.WorkerPatch(ctx, orgmap.UpdateWorkerToModel(worker)); err != nil {
		o.Logger.Error(
			"failed to update worker",
			zap.Error(err),
		)
		return err
	}
	return nil
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

func (o *OrgUseCase) WorkerList(ctx context.Context, orgID, limit, page int) (*orgdto.WorkerList, error) {
	offset := (page - 1) * limit
	data, found, err := o.org.WorkerList(ctx, orgID, limit, offset)
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
	resp := &orgdto.WorkerList{
		List:  workerList,
		Found: found,
	}
	return resp, nil
}
func (o *OrgUseCase) WorkerDelete(ctx context.Context, workerID, orgID int) error {
	if err := o.org.WorkerSoftDelete(ctx, workerID, orgID); err != nil {
		o.Logger.Error(
			"failed to delete worker",
			zap.Error(err),
		)
		return err
	}
	return nil
}
