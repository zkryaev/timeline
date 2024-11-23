package facemap

import (
	"timeline/internal/entity/dto"
	"timeline/internal/repository/models"
)

func SearchToModel(dto *dto.SearchReq) *models.SearchParams {
	return &models.SearchParams{
		Page:   dto.Page,
		Limit:  dto.Limit,
		Offset: (dto.Page - 1) * dto.Limit,
		Name:   dto.Name,
		Type:   dto.Type,
	}
}

func AreaToModel(dto *dto.OrgAreaReq) *models.AreaParams {
	return &models.AreaParams{
		Left:  models.Coordinate{Lat: dto.LeftLowerCorner.Lat, Long: dto.LeftLowerCorner.Long},
		Right: models.Coordinate{Lat: dto.LeftLowerCorner.Lat, Long: dto.LeftLowerCorner.Long},
	}
}
