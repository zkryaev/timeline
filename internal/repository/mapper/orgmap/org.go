package orgmap

import (
	"database/sql"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
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

func OrgInfoToDTO(model *models.OrgInfo) *orgdto.Organization {
	return &orgdto.Organization{
		OrgID: model.OrgID,
		Info:  OrgInfoToEntity(model),
	}
}

func OrgInfoToEntity(model *models.OrgInfo) *entity.OrgInfo {
	resp := &entity.OrgInfo{
		Name:      model.Name,
		Rating:    model.Rating,
		Address:   model.Address,
		Long:      model.Long,
		Lat:       model.Lat,
		Type:      model.Type,
		Telephone: model.Telephone,
		City:      model.City,
		About:     model.About,
		Timetable: TimetableToEntity(model.Timetable),
	}
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

func MapOrgInfoToModel(dto *entity.MapOrgInfo) *models.OrgSummary {
	return &models.OrgSummary{
		OrgID:      dto.OrgID,
		Name:       dto.Name,
		Rating:     dto.Rating,
		Type:       dto.Type,
		OpenHours:  *OpenHoursToModel(dto.TodaySchedule),
		Coordinate: models.Coordinate{Lat: dto.Coords.Lat, Long: dto.Coords.Long},
	}
}

func OrgSummaryToDTO(model *models.OrgSummary) *entity.MapOrgInfo {
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
		Timetable: TimetableToModel(dto.Timetable),
	}
}

// func UpdateToDTO(model *models.OrgUpdate) *orgdto.OrgUpdateResp {
// 	return &orgdto.OrgUpdateResp{
// 		OrgID:     model.OrgID,
// 		Name:      model.Name,
// 		Type:      model.Type,
// 		City:      model.City,
// 		Address:   model.Address,
// 		Telephone: model.Telephone,
// 		Lat:       model.Lat,
// 		Long:      model.Long,
// 		About:     model.About,
// 		Timetable: TimetableToEntity(model.Timetable),
// 	}
// }
