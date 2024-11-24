package usermap

import (
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/entity/dto/userdto"
	"timeline/internal/repository/models"
)

func UserRegisterToModel(dto *authdto.UserRegisterReq) *models.UserRegister {
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

func UserUpdateToModel(dto *userdto.UserUpdateReq) *models.UserInfo {
	return &models.UserInfo{
		UserID:    dto.UserID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Telephone: dto.Telephone,
		City:      dto.City,
		About:     dto.About,
	}
}

func UserUpdateToDTO(model *models.UserInfo) *userdto.UserUpdateResp {
	return &userdto.UserUpdateResp{
		UserID:    model.UserID,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Telephone: model.Telephone,
		City:      model.City,
		About:     model.About,
	}
}
