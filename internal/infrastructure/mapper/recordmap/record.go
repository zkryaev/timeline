package recordmap

import (
	"time"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/mapper/usermap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/orgmodel"
	"timeline/internal/infrastructure/models/recordmodel"
)

func RecordToModel(dto *recordto.Record) *recordmodel.Record {
	return &recordmodel.Record{
		RecordID:  dto.RecordID,
		OrgID:     dto.OrgID,
		UserID:    dto.UserID,
		SlotID:    dto.SlotID,
		ServiceID: dto.ServiceID,
		WorkerID:  dto.WorkerID,
		Reviewed:  dto.Reviewed,
	}
}

func RecordToDTO(model *recordmodel.Record) *recordto.Record {
	return &recordto.Record{
		RecordID:  model.RecordID,
		OrgID:     model.OrgID,
		UserID:    model.UserID,
		SlotID:    model.SlotID,
		ServiceID: model.ServiceID,
		WorkerID:  model.WorkerID,
		Reviewed:  model.Reviewed,
	}
}

func RecordParamToModel(dto recordto.RecordParam) recordmodel.RecordParam {
	return recordmodel.RecordParam{
		RecordID: dto.RecordID,
		TData:    models.TokenData(dto.TData),
	}
}

func RecordParamsToModel(dto *recordto.RecordListParams) *recordmodel.RecordListParams {
	return &recordmodel.RecordListParams{
		OrgID:    dto.OrgID,
		UserID:   dto.UserID,
		Reviewed: dto.Reviewed,
		Fresh:    dto.Fresh,
		Limit:    dto.Limit,
		Offset:   (dto.Page - 1) * dto.Limit,
		TData:    models.TokenData(dto.TData),
	}
}

func RecordScrapToDTO(model *recordmodel.RecordScrap, loc *time.Location) *recordto.RecordScrap {
	return &recordto.RecordScrap{
		RecordID: model.RecordID,
		Reviewed: model.Reviewed,
		Org:      orgmap.OrganizationToDTO(&orgmodel.Organization{OrgInfo: *model.Org}, loc),
		User:     usermap.UserInfoToDTO(model.User),
		Slot:     orgmap.SlotInfoToDTO(model.Slot, loc),
		Service:  orgmap.ServiceToEntity(model.Service),
		Worker:   orgmap.WorkerToEntity(model.Worker),
		Feedback: FeedbackToDTO(model.Feedback),
	}
}

func RecordListToDTO(model []*recordmodel.RecordScrap, loc *time.Location) []*recordto.RecordScrap {
	list := make([]*recordto.RecordScrap, 0, len(model))
	for _, v := range model {
		list = append(list, RecordScrapToDTO(v, loc))
	}
	return list
}

func ReminderRecordToReminder(model *recordmodel.ReminderRecord, loc *time.Location) *models.ReminderMsg {
	return &models.ReminderMsg{
		Organization: model.OrgName,
		Service:      model.ServiceName,
		ServiceDesc:  model.ServiceDescription,
		Address:      model.OrgAddress,
		SessionStart: time.Date(model.Date.Year(), model.Date.Month(), model.Date.Day(),
			model.Begin.Hour(), model.Begin.Minute(), model.Begin.Second(),
			model.Begin.Nanosecond(), model.Date.Location(),
		).In(loc),
		SessionEnd: time.Date(model.Date.Year(), model.Date.Month(), model.Date.Day(),
			model.End.Hour(), model.End.Minute(), model.End.Second(),
			model.End.Nanosecond(), model.Date.Location(),
		).In(loc),
		SessionDate: model.Date.In(loc),
	}
}

func CancelationToModel(dto *recordto.RecordCancelation) *recordmodel.RecordCancelation {
	return &recordmodel.RecordCancelation{
		TData:        models.TokenData(dto.TData),
		CancelReason: dto.CancelReason,
		RecordID:     dto.RecordID,
		IsCanceled:   true,
	}
}
