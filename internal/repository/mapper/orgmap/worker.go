package orgmap

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/models"
)

func AddWorkerToModel(dto *orgdto.AddWorkerReq) *models.Worker {
	return &models.Worker{
		WorkerID:  0,
		OrgID:     dto.OrgID,
		FirstName: dto.WorkerInfo.FirstName,
		LastName:  dto.WorkerInfo.LastName,
		Position:  dto.WorkerInfo.Position,
		Degree:    dto.WorkerInfo.Degree,
	}
}

func UpdateWorkerToModel(dto *orgdto.UpdateWorkerReq) *models.Worker {
	return &models.Worker{
		WorkerID:  dto.WorkerID,
		OrgID:     dto.OrgID,
		FirstName: dto.WorkerInfo.FirstName,
		LastName:  dto.WorkerInfo.LastName,
		Position:  dto.WorkerInfo.Position,
		Degree:    dto.WorkerInfo.Degree,
	}
}

func WorkerToDTO(model *models.Worker) *orgdto.WorkerResp {
	return &orgdto.WorkerResp{
		WorkerID:   model.WorkerID,
		OrgID:      model.OrgID,
		WorkerInfo: workerToEntity(model),
	}
}

func workerToEntity(model *models.Worker) *entity.Worker {
	return &entity.Worker{
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Position:  model.Position,
		Degree:    model.Degree,
	}
}

func AssignWorkerToModel(dto *orgdto.AssignWorkerReq) *models.WorkerAssign {
	return &models.WorkerAssign{
		WorkerID:  dto.WorkerID,
		OrgID:     dto.OrgID,
		ServiceID: dto.ServiceID,
	}
}
