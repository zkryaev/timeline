package orgmap

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/models"
)

func RegisterReqToModel(dto *authdto.OrgRegisterReq) *models.OrgRegister {
	return &models.OrgRegister{
		HashCreds: models.HashCreds{
			Email:      dto.Email,
			PasswdHash: dto.Password,
		},
		OrgInfo: models.OrgInfo{
			Name:      dto.Name,
			Rating:    dto.Rating,
			Type:      dto.Type,
			City:      dto.City,
			Address:   dto.Address,
			Telephone: dto.Telephone,
			Long:      dto.Long,
			Lat:       dto.Lat,
			About:     dto.About,
		},
	}
}

func OrgInfoToDTO(model *models.OrgInfo) *entity.Organization {
	return &entity.Organization{
		OrgID: model.OrgID,
		Info:  *OrgInfoToEntity(model),
	}
}

func OrgInfoToEntity(model *models.OrgInfo) *entity.OrgInfo {
	return &entity.OrgInfo{
		Name:      model.Name,
		Rating:    model.Rating,
		Address:   model.Address,
		Long:      model.Long,
		Lat:       model.Lat,
		Type:      model.Type,
		Telephone: model.Telephone,
		City:      model.City,
		About:     model.About,
	}
}

func MapOrgInfoToModel(dto *entity.MapOrgInfo) *models.OrgSummary {
	return &models.OrgSummary{
		OrgID:      dto.OrgID,
		Name:       dto.Name,
		Rating:     dto.Rating,
		Type:       dto.Type,
		Coordinate: models.Coordinate{Lat: dto.Coords.Lat, Long: dto.Coords.Long},
	}
}

func OrgSummaryToDTO(model *models.OrgSummary) *entity.MapOrgInfo {
	return &entity.MapOrgInfo{
		OrgID:  model.OrgID,
		Name:   model.Name,
		Rating: model.Rating,
		Type:   model.Type,
		Coords: entity.MapPoint{Lat: model.Coordinate.Lat, Long: model.Coordinate.Long},
	}
}

func UpdateToModel(dto *orgdto.OrgUpdateReq) *models.OrgUpdate {
	return &models.OrgUpdate{
		OrgID:     dto.OrgID,
		Name:      dto.Name,
		Type:      dto.Type,
		City:      dto.City,
		Address:   dto.Address,
		Telephone: dto.Telephone,
		Lat:       dto.Lat,
		Long:      dto.Long,
		About:     dto.About,
	}
}

func UpdateToDTO(model *models.OrgUpdate) *orgdto.OrgUpdateResp {
	return &orgdto.OrgUpdateResp{
		OrgID:     model.OrgID,
		Name:      model.Name,
		Type:      model.Type,
		City:      model.City,
		Address:   model.Address,
		Telephone: model.Telephone,
		Lat:       model.Lat,
		Long:      model.Long,
		About:     model.About,
	}
}
