package orgmap

import (
	"strings"
	"time"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/models/orgmodel"
)

func SlotToDTO(model *orgmodel.Slot, loc *time.Location) *orgdto.SlotResp {
	return &orgdto.SlotResp{
		SlotID: model.SlotID,
		Slot:   *SlotInfoToDTO(model, loc),
	}
}

func SlotReqToModel(dto *orgdto.SlotReq) *orgmodel.SlotsMeta {
	return &orgmodel.SlotsMeta{
		SlotID:   dto.SlotID,
		WorkerID: dto.WorkerID,
		UserID:   dto.UserID,
		OrgID:    dto.OrgID,
	}
}

func SlotUpdateToModel(dto *orgdto.SlotUpdate) *orgmodel.SlotsMeta {
	return &orgmodel.SlotsMeta{
		SlotID:   dto.SlotID,
		WorkerID: dto.WorkerID,
	}
}

func SlotInfoToDTO(model *orgmodel.Slot, loc *time.Location) *orgdto.Slot {
	return &orgdto.Slot{
		WorkerScheduleID: model.WorkerScheduleID,
		WorkerID:         model.WorkerID,
		Date:             strings.Fields(model.Date.String())[0],
		Begin:            model.Begin.In(loc).Format(timeFormat),
		End:              model.End.In(loc).Format(timeFormat),
		Busy:             model.Busy,
	}
}
