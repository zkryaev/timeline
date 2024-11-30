package orgmap

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/models"
)

func ServiceToDTO(model *models.Service) *orgdto.ServiceResp {
	return &orgdto.ServiceResp{
		ServiceID:   model.ServiceID,
		OrgID:       model.OrgID,
		ServiceInfo: serviceToEntity(model),
	}
}

func AddServiceToModel(dto *orgdto.AddServiceReq) *models.Service {
	return &models.Service{
		ServiceID:   0, // zeroval
		OrgID:       dto.OrgID,
		Name:        dto.ServiceInfo.Name,
		Cost:        dto.ServiceInfo.Cost,
		Description: dto.ServiceInfo.Description,
	}
}

func UpdateService(dto *orgdto.UpdateServiceReq) *models.Service {
	return &models.Service{
		ServiceID:   dto.ServiceID,
		OrgID:       dto.OrgID,
		Name:        dto.ServiceInfo.Name,
		Cost:        dto.ServiceInfo.Cost,
		Description: dto.ServiceInfo.Description,
	}
}

func serviceToEntity(model *models.Service) *entity.Service {
	return &entity.Service{
		Name:        model.Name,
		Cost:        model.Cost,
		Description: model.Description,
	}
}
