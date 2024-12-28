package orgmap

import (
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/mediamap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/orgmodel"
)

func RegisterReqToModel(dto *authdto.OrgRegisterReq) *orgmodel.OrgRegister {
	return &orgmodel.OrgRegister{
		UUID: dto.UUID,
		HashCreds: models.HashCreds{
			Email:      dto.Email,
			PasswdHash: dto.Password,
		},
		OrgInfo: orgmodel.OrgInfo{
			Name:        dto.Name,
			Rating:      dto.Rating,
			Type:        dto.Type,
			City:        dto.City,
			Address:     dto.Address,
			Telephone:   dto.Telephone,
			Coordinates: *CoordsToModel(&dto.Coordinates),
			About:       dto.About,
		},
	}
}

func CoordsToEntity(model *orgmodel.Coordinates) *entity.Coordinates {
	return &entity.Coordinates{
		Lat:  model.Lat,
		Long: model.Long,
	}
}

func CoordsToModel(dto *entity.Coordinates) *orgmodel.Coordinates {
	return &orgmodel.Coordinates{
		Lat:  dto.Lat,
		Long: dto.Long,
	}
}

func OrganizationToDTO(model *orgmodel.Organization) *orgdto.Organization {
	return &orgdto.Organization{
		ImagesURL: mediamap.ImageUUIDToURL(model.ImagesURL...),
		OrgID:     model.OrgID,
		Info:      OrgInfoToEntity(&model.OrgInfo),
		Timetable: TimetableToEntity(model.Timetable),
	}
}

func OrgsBySearchToDTO(model *orgmodel.OrgsBySearch) *entity.OrgsBySearch {
	return &entity.OrgsBySearch{
		OrgID:         model.OrgID,
		Name:          model.Name,
		Rating:        model.Rating,
		Type:          model.Type,
		Address:       model.Address,
		Coords:        CoordsToEntity(&model.Coordinates),
		TodaySchedule: OpenHoursToDTO(&model.OpenHours),
	}
}

func OrgInfoToEntity(model *orgmodel.OrgInfo) *entity.OrgInfo {
	resp := &entity.OrgInfo{
		Name:        model.Name,
		Rating:      model.Rating,
		Address:     model.Address,
		Coordinates: *CoordsToEntity(&model.Coordinates),
		Type:        model.Type,
		Telephone:   model.Telephone,
		City:        model.City,
		About:       model.About,
	}
	return resp
}

func SearchToModel(dto *general.SearchReq) *orgmodel.SearchParams {
	return &orgmodel.SearchParams{
		Page:   dto.Page,
		Limit:  dto.Limit,
		Offset: (dto.Page - 1) * dto.Limit,
		Name:   dto.Name,
		Type:   dto.Type,
	}
}

func AreaToModel(dto *general.OrgAreaReq) *orgmodel.AreaParams {
	return &orgmodel.AreaParams{
		Left:  *CoordsToModel(&dto.LeftLowerCorner),
		Right: *CoordsToModel(&dto.RightUpperCorner),
	}
}

func OrgInfoToModel(model *entity.OrgInfo) *orgmodel.OrgInfo {
	resp := &orgmodel.OrgInfo{
		Name:        model.Name,
		Rating:      model.Rating,
		Address:     model.Address,
		Coordinates: *CoordsToModel(&model.Coordinates),
		Type:        model.Type,
		Telephone:   model.Telephone,
		City:        model.City,
		About:       model.About,
	}
	return resp
}

func OrgUpdateToModel(dto *orgdto.OrgUpdateReq) *orgmodel.Organization {
	resp := &orgmodel.Organization{}
	resp.OrgID = dto.OrgID
	resp.OrgInfo = *OrgInfoToModel(&dto.OrgInfo)
	resp.Timetable = TimetableToModel(dto.Timetable)
	return resp
}

func TimetableToEntity(Timetable []*orgmodel.OpenHours) []*entity.OpenHours {
	if len(Timetable) == 0 {
		return nil
	}
	resp := make([]*entity.OpenHours, 0, len(Timetable))
	for _, v := range Timetable {
		resp = append(resp, OpenHoursToDTO(v))
	}
	return resp
}

func TimetableToModel(Timetable []*entity.OpenHours) []*orgmodel.OpenHours {
	if len(Timetable) == 0 {
		return nil
	}
	resp := make([]*orgmodel.OpenHours, 0, 2)
	for _, v := range Timetable {
		resp = append(resp, OpenHoursToModel(v))
	}
	return resp
}

func MapOrgInfoToModel(dto *entity.MapOrgInfo) *orgmodel.OrgByArea {
	return &orgmodel.OrgByArea{
		OrgID:       dto.OrgID,
		Name:        dto.Name,
		Rating:      dto.Rating,
		Type:        dto.Type,
		OpenHours:   *OpenHoursToModel(dto.TodaySchedule),
		Coordinates: *CoordsToModel(&dto.Coords),
	}
}

func OrgSummaryToDTO(model *orgmodel.OrgByArea) *entity.MapOrgInfo {
	return &entity.MapOrgInfo{
		OrgID:  model.OrgID,
		Name:   model.Name,
		Rating: model.Rating,
		Type:   model.Type,
		TodaySchedule: OpenHoursToDTO(
			&orgmodel.OpenHours{
				Weekday:    model.Weekday,
				Open:       model.Open,
				Close:      model.Close,
				BreakStart: model.BreakStart,
				BreakEnd:   model.BreakEnd,
			},
		),
		Coords: *CoordsToEntity(&model.Coordinates),
	}
}
