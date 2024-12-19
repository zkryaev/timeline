package recordmap

import (
	"time"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/repository/mail/entity"
	"timeline/internal/repository/mapper/orgmap"
	"timeline/internal/repository/mapper/usermap"
	"timeline/internal/repository/models/recordmodel"
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

func RecordParamsToModel(dto *recordto.RecordListParams) *recordmodel.RecordListParams {
	return &recordmodel.RecordListParams{
		OrgID:    dto.OrgID,
		UserID:   dto.UserID,
		Reviewed: dto.Reviewed,
		Fresh:    dto.Fresh,
	}
}

func RecordScrapToDTO(model *recordmodel.RecordScrap) *recordto.RecordScrap {
	return &recordto.RecordScrap{
		RecordID: model.RecordID,
		Reviewed: model.Reviewed,
		Org:      orgmap.OrgInfoToEntity(model.Org),
		User:     usermap.UserInfoToDTO(model.User),
		Slot:     orgmap.SlotInfoToDTO(model.Slot),
		Service:  orgmap.ServiceToEntity(model.Service),
		Worker:   orgmap.WorkerToEntity(model.Worker),
		Feedback: FeedbackToDTO(model.Feedback),
	}
}

func RecordListToDTO(model []*recordmodel.RecordScrap) []*recordto.RecordScrap {
	list := make([]*recordto.RecordScrap, 0, len(model))
	for _, v := range model {
		list = append(list, RecordScrapToDTO(v))
	}
	return list
}

func RecordToReminder(model *recordmodel.ReminderRecord) *entity.ReminderMsg {
	return &entity.ReminderMsg{
		Organization: model.OrgName,
		Service:      model.ServiceName,
		Description:  model.ServiceDescription,
		Address:      model.OrgAddress,
		SessionStart: time.Date(model.Date.Year(), model.Date.Month(), model.Date.Day(), model.Begin.Hour(), model.Begin.Minute(), model.Begin.Second(), model.Begin.Nanosecond(), model.Begin.UTC().Location()),
		SessionEnd:   time.Date(model.Date.Year(), model.Date.Month(), model.Date.Day(), model.End.Hour(), model.End.Minute(), model.End.Second(), model.End.Nanosecond(), model.End.UTC().Location()),
		SessionDate:  model.Date,
	}
}
