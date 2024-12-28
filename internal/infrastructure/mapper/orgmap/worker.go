package orgmap

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/models/orgmodel"
)

func AddWorkerToModel(dto *orgdto.AddWorkerReq) *orgmodel.Worker {
	return &orgmodel.Worker{
		WorkerID:        0,
		OrgID:           dto.OrgID,
		FirstName:       dto.WorkerInfo.FirstName,
		LastName:        dto.WorkerInfo.LastName,
		Position:        dto.WorkerInfo.Position,
		Degree:          dto.WorkerInfo.Degree,
		SessionDuration: dto.WorkerInfo.SessionDuration,
	}
}

func UpdateWorkerToModel(dto *orgdto.UpdateWorkerReq) *orgmodel.Worker {
	return &orgmodel.Worker{
		WorkerID:        dto.WorkerID,
		OrgID:           dto.OrgID,
		FirstName:       dto.WorkerInfo.FirstName,
		LastName:        dto.WorkerInfo.LastName,
		Position:        dto.WorkerInfo.Position,
		Degree:          dto.WorkerInfo.Degree,
		SessionDuration: dto.WorkerInfo.SessionDuration,
	}
}

func WorkerToDTO(model *orgmodel.Worker) *orgdto.WorkerResp {
	return &orgdto.WorkerResp{
		WorkerID:   model.WorkerID,
		OrgID:      model.OrgID,
		WorkerInfo: WorkerToEntity(model),
	}
}

func WorkerToEntity(model *orgmodel.Worker) *entity.Worker {
	return &entity.Worker{
		UUID:            model.UUID,
		FirstName:       model.FirstName,
		LastName:        model.LastName,
		Position:        model.Position,
		Degree:          model.Degree,
		SessionDuration: model.SessionDuration,
	}
}

func AssignWorkerToModel(dto *orgdto.AssignWorkerReq) *orgmodel.WorkerAssign {
	return &orgmodel.WorkerAssign{
		WorkerID:  dto.WorkerID,
		OrgID:     dto.OrgID,
		ServiceID: dto.ServiceID,
	}
}
