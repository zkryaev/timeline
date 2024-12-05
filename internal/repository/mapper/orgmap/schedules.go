package orgmap

import (
	"time"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/models/orgmodel"
)

func ScheduleListToModel(dto *orgdto.ScheduleList) *orgmodel.ScheduleList {
	resp := &orgmodel.ScheduleList{
		WorkerID: dto.WorkerID,
		OrgID:    dto.OrgID,
		Schedule: make([]*orgmodel.Schedule, 0, len(dto.Schedule)),
	}
	for _, v := range dto.Schedule {
		resp.Schedule = append(resp.Schedule, scheduleToModel(v))
	}
	return resp
}

func ScheduleListToDTO(model *orgmodel.ScheduleList) *orgdto.ScheduleList {
	resp := &orgdto.ScheduleList{
		WorkerID: model.WorkerID,
		OrgID:    model.OrgID,
		Schedule: make([]*orgdto.Schedule, 0, len(model.Schedule)),
	}
	for _, v := range model.Schedule {
		resp.Schedule = append(resp.Schedule, scheduleToDTO(v))
	}
	return resp
}

func ScheduleParamsToModel(dto *orgdto.ScheduleParams) *orgmodel.ScheduleParams {
	return &orgmodel.ScheduleParams{
		WorkerID: dto.WorkerID,
		OrgID:    dto.OrgID,
		Weekday:  dto.Weekday,
	}
}

func scheduleToModel(dto *orgdto.Schedule) *orgmodel.Schedule {
	start, _ := time.Parse(timeFormat, dto.Start)
	over, _ := time.Parse(timeFormat, dto.Over)
	start = start.AddDate(2001, 0, 0) // Прибавляем 2001 год
	over = over.AddDate(2001, 0, 0)   // Прибавляем 2001 год
	return &orgmodel.Schedule{
		WorkerScheduleID: dto.WorkerScheduleID,
		Weekday:          dto.Weekday,
		Start:            start,
		Over:             over,
	}
}

func scheduleToDTO(model *orgmodel.Schedule) *orgdto.Schedule {
	return &orgdto.Schedule{
		WorkerScheduleID: model.WorkerScheduleID,
		Weekday:          model.Weekday,
		Start:            model.Start.Format(timeFormat),
		Over:             model.Over.Format(timeFormat),
	}
}
