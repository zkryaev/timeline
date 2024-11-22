package orgmap

import (
	"timeline/internal/entity/dto"
	"timeline/internal/repository/models"
)

func ToModel(dto *dto.OrgRegisterReq) *models.OrgRegister {
	return &models.OrgRegister{
		HashCreds: models.HashCreds{
			Email:      dto.Email,
			PasswdHash: dto.Password,
		},
		OrgInfo: models.OrgInfo{
			Name:      dto.Name,
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

func ToDTO(model *models.OrgInfo) {}
