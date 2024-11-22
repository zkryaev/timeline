package usermap

import (
	"timeline/internal/entity/dto"
	"timeline/internal/repository/models"
)

func ToModel(dto *dto.UserRegisterReq) *models.UserRegister {
	return &models.UserRegister{
		HashCreds: models.HashCreds{
			Email:      dto.Email,
			PasswdHash: dto.Password,
		},
		UserInfo: models.UserInfo{
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Telephone: dto.Telephone,
			City:      dto.City,
			About:     dto.About,
		},
	}
}
