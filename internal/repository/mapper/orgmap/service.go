package orgmap

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/models/orgmodel"
)

func ServiceToDTO(model *orgmodel.Service) *orgdto.ServiceResp {
	return &orgdto.ServiceResp{
		ServiceID:   model.ServiceID,
		OrgID:       model.OrgID,
		ServiceInfo: ServiceToEntity(model),
	}
}

func AddServiceToModel(dto *orgdto.AddServiceReq) *orgmodel.Service {
	return &orgmodel.Service{
		ServiceID:   0, // zeroval
		OrgID:       dto.OrgID,
		Name:        dto.ServiceInfo.Name,
		Cost:        dto.ServiceInfo.Cost,
		Description: dto.ServiceInfo.Description,
	}
}

func UpdateService(dto *orgdto.UpdateServiceReq) *orgmodel.Service {
	return &orgmodel.Service{
		ServiceID:   dto.ServiceID,
		OrgID:       dto.OrgID,
		Name:        dto.ServiceInfo.Name,
		Cost:        dto.ServiceInfo.Cost,
		Description: dto.ServiceInfo.Description,
	}
}

func ServiceToEntity(model *orgmodel.Service) *entity.Service {
	return &entity.Service{
		Name:        model.Name,
		Cost:        model.Cost,
		Description: model.Description,
	}
}
