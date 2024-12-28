package orgmap

import (
	"time"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/models/orgmodel"
)

func WorkerScheduleToModel(dto *orgdto.WorkerSchedule) *orgmodel.WorkerSchedule {
	resp := &orgmodel.WorkerSchedule{
		WorkerID:        dto.WorkerID,
		OrgID:           dto.OrgID,
		SessionDuration: dto.SessionDuration,
		Schedule:        make([]*orgmodel.Schedule, 0, len(dto.Schedule)),
	}
	for _, v := range dto.Schedule {
		resp.Schedule = append(resp.Schedule, scheduleToModel(v))
	}
	return resp
}

func WorkerScheduleToDTO(model *orgmodel.WorkerSchedule) *orgdto.WorkerSchedule {
	resp := &orgdto.WorkerSchedule{
		WorkerID:        model.WorkerID,
		OrgID:           model.OrgID,
		SessionDuration: model.SessionDuration,
		Schedule:        make([]*orgdto.Schedule, 0, len(model.Schedule)),
	}
	for _, v := range model.Schedule {
		resp.Schedule = append(resp.Schedule, scheduleToDTO(v))
	}
	return resp
}

func ScheduleListToDTO(model *orgmodel.ScheduleList) *orgdto.ScheduleList {
	resp := &orgdto.ScheduleList{
		Workers: make([]*orgdto.WorkerSchedule, 0, len(model.Workers)),
		Found:   model.Found,
	}
	for _, v := range model.Workers {
		resp.Workers = append(resp.Workers, WorkerScheduleToDTO(v))
	}
	return resp
}

func ScheduleParamsToModel(dto *orgdto.ScheduleParams) *orgmodel.ScheduleParams {
	return &orgmodel.ScheduleParams{
		WorkerID: dto.WorkerID,
		OrgID:    dto.OrgID,
		Weekday:  dto.Weekday,
		Limit:    dto.Limit,
		Offset:   (dto.Page - 1) * dto.Limit,
	}
}

func scheduleToModel(dto *orgdto.Schedule) *orgmodel.Schedule {
	start, _ := time.Parse(timeFormat, dto.Start)
	over, _ := time.Parse(timeFormat, dto.Over)
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
