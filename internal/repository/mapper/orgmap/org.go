package orgmap

import (
	"database/sql"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/models"
)

const timeFormat = "15:04"

func RegisterReqToModel(dto *authdto.OrgRegisterReq) *models.OrgRegister {
	return &models.OrgRegister{
		HashCreds: models.HashCreds{
			Email:      dto.Email,
			PasswdHash: dto.Password,
		},
		OrgInfo: models.OrgInfo{
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

func CoordsToEntity(model *models.Coordinates) *entity.Coordinates {
	return &entity.Coordinates{
		Lat:  model.Lat,
		Long: model.Long,
	}
}

func CoordsToModel(dto *entity.Coordinates) *models.Coordinates {
	return &models.Coordinates{
		Lat:  dto.Lat,
		Long: dto.Long,
	}
}

func OrganizationToDTO(model *models.Organization) *orgdto.Organization {
	return &orgdto.Organization{
		OrgID:     model.OrgID,
		Info:      OrgInfoToEntity(&model.OrgInfo),
		Timetable: TimetableToEntity(model.Timetable),
	}
}

func OrgsBySearchToDTO(model *models.OrgsBySearch) *entity.OrgsBySearch {
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

func OrgInfoToEntity(model *models.OrgInfo) *entity.OrgInfo {
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

func SearchToModel(dto *general.SearchReq) *models.SearchParams {
	return &models.SearchParams{
		Page:   dto.Page,
		Limit:  dto.Limit,
		Offset: (dto.Page - 1) * dto.Limit,
		Name:   dto.Name,
		Type:   dto.Type,
	}
}

func AreaToModel(dto *general.OrgAreaReq) *models.AreaParams {
	return &models.AreaParams{
		Left:  *CoordsToModel(&dto.LeftLowerCorner),
		Right: *CoordsToModel(&dto.RightUpperCorner),
	}
}

func OrgInfoToModel(model *entity.OrgInfo) *models.OrgInfo {
	resp := &models.OrgInfo{
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

func OrgUpdateToModel(dto *orgdto.OrgUpdateReq) *models.Organization {
	resp := &models.Organization{}
	resp.OrgID = dto.OrgID
	resp.OrgInfo = *OrgInfoToModel(&dto.OrgInfo)
	resp.Timetable = TimetableToModel(dto.Timetable)
	return resp
}

func TimetableToEntity(Timetable []*models.OpenHours) []*entity.OpenHours {
	if len(Timetable) == 0 {
		return nil
	}
	resp := make([]*entity.OpenHours, 0, len(Timetable))
	for _, v := range Timetable {
		resp = append(resp, OpenHoursToDTO(v))
	}
	return resp
}

func TimetableToModel(Timetable []*entity.OpenHours) []*models.OpenHours {
	if len(Timetable) == 0 {
		return nil
	}
	resp := make([]*models.OpenHours, 0, 2)
	for _, v := range Timetable {
		resp = append(resp, OpenHoursToModel(v))
	}
	return resp
}

func OpenHoursToModel(day *entity.OpenHours) *models.OpenHours {
	open, _ := time.Parse(timeFormat, day.Open)
	close, _ := time.Parse(timeFormat, day.Close)
	breakstart, _ := time.Parse(timeFormat, day.BreakStart)
	breakend, _ := time.Parse(timeFormat, day.BreakEnd)
	return &models.OpenHours{
		Weekday:    sql.NullInt32{Int32: int32(day.Weekday), Valid: true},
		Open:       sql.NullTime{Time: open, Valid: true},
		Close:      sql.NullTime{Time: close, Valid: true},
		BreakStart: sql.NullTime{Time: breakstart, Valid: true},
		BreakEnd:   sql.NullTime{Time: breakend, Valid: true},
	}
}

func OpenHoursToDTO(day *models.OpenHours) *entity.OpenHours {
	// Если weekday = -1 значит пустая структура
	if !day.Weekday.Valid {
		return nil
	}
	return &entity.OpenHours{
		Weekday:    int(day.Weekday.Int32),
		Open:       day.Open.Time.Format(timeFormat),
		Close:      day.Close.Time.Format(timeFormat),
		BreakStart: day.BreakStart.Time.Format(timeFormat),
		BreakEnd:   day.BreakEnd.Time.Format(timeFormat),
	}
}

func MapOrgInfoToModel(dto *entity.MapOrgInfo) *models.OrgByArea {
	return &models.OrgByArea{
		OrgID:       dto.OrgID,
		Name:        dto.Name,
		Rating:      dto.Rating,
		Type:        dto.Type,
		OpenHours:   *OpenHoursToModel(dto.TodaySchedule),
		Coordinates: *CoordsToModel(&dto.Coords),
	}
}

func OrgSummaryToDTO(model *models.OrgByArea) *entity.MapOrgInfo {
	return &entity.MapOrgInfo{
		OrgID:  model.OrgID,
		Name:   model.Name,
		Rating: model.Rating,
		Type:   model.Type,
		TodaySchedule: OpenHoursToDTO(
			&models.OpenHours{
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
