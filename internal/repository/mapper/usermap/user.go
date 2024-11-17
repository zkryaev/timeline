package usermap

import (
	"timeline/internal/entity/dto"
	"timeline/internal/repository/models"
)

func ToModel(dto *dto.UserRegisterReq) *models.UserRegisterModel {
	return &models.UserRegisterModel{
		HashCreds: models.HashCreds{
			Email:      dto.Email,
			PasswdHash: dto.Password,
		},
		UserInfo: models.UserInfo{
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Telephone: dto.Telephone,
			Social:    dto.Social,
			About:     dto.About,
		},
	}
}
