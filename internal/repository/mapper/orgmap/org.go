package orgmap

import (
	"timeline/internal/entity/dto"
	"timeline/internal/repository/models"
)

func ToModel(dto *dto.OrgRegisterReq) *models.OrgRegisterModel {
	return &models.OrgRegisterModel{
		HashCreds: models.HashCreds{
			Email:      dto.Email,
			PasswdHash: dto.Password,
		},
		OrgInfo: models.OrgInfo{
			Name:      dto.Name,
			Address:   dto.Address,
			Long:      dto.Long,
			Lat:       dto.Lat,
			Type:      dto.Type,
			Telephone: dto.Telephone,
			Social:    dto.Social,
			About:     dto.About,
		},
	}
}

func ToDTO(model *models.OrgInfo) {}
