package orgmap

import (
	"time"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/orgmodel"
	"timeline/internal/usecase/common"
)

func SlotToDTO(model *orgmodel.Slot, loc *time.Location) *orgdto.SlotResp {
	return &orgdto.SlotResp{
		SlotID: model.SlotID,
		Slot:   *SlotInfoToDTO(model, loc),
	}
}

func SlotReqToModel(dto *orgdto.SlotReq) *orgmodel.SlotsReq {
	return &orgmodel.SlotsReq{
		WorkerID: dto.WorkerID,
		OrgID:    dto.OrgID,
		TData:    models.TokenData(dto.TData),
	}
}

func SlotInfoToDTO(model *orgmodel.Slot, loc *time.Location) *orgdto.Slot {
	return &orgdto.Slot{
		WorkerScheduleID: model.WorkerScheduleID,
		WorkerID:         model.WorkerID,
		Date:             model.Date.Format(time.DateOnly),
		Begin:            model.Begin.In(loc).Format(common.MinutesOnly),
		End:              model.End.In(loc).Format(common.MinutesOnly),
		Busy:             model.Busy,
	}
}
