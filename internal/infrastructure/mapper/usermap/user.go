package usermap

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/entity/dto/userdto"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/usermodel"
)

func UserRegisterToModel(dto *authdto.UserRegisterReq) *usermodel.UserRegister {
	return &usermodel.UserRegister{
		UUID: dto.UUID,
		HashCreds: models.HashCreds{
			Email:      dto.Credentials.Email,
			PasswdHash: dto.Credentials.Password,
		},
		UserInfo: usermodel.UserInfo{
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Telephone: dto.Telephone,
			City:      dto.City,
			About:     dto.About,
		},
	}
}

func UserUpdateToModel(dto *userdto.UserUpdateReq) *usermodel.UserInfo {
	return &usermodel.UserInfo{
		UserID:    dto.UserID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Telephone: dto.Telephone,
		City:      dto.City,
		About:     dto.About,
	}
}

func UserInfoToEntity(dto *entity.User) *usermodel.UserInfo {
	return &usermodel.UserInfo{
		UserID:    dto.UserID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Telephone: dto.Telephone,
		City:      dto.City,
		About:     dto.About,
	}
}

func UserInfoToDTO(model *usermodel.UserInfo) *entity.User {
	return &entity.User{
		UserID:    model.UserID,
		UUID:      model.UUID,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Telephone: model.Telephone,
		City:      model.City,
		About:     model.About,
		Email:     model.Email,
	}
}
