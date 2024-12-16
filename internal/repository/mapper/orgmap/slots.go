package orgmap

import (
	"strings"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/models/orgmodel"
)

func SlotToDTO(model *orgmodel.Slot) *orgdto.SlotResp {
	return &orgdto.SlotResp{
		SlotID: model.SlotID,
		Slot:   *SlotInfoToDTO(model),
	}
}

func SlotReqToModel(dto *orgdto.SlotReq) *orgmodel.SlotsMeta {
	return &orgmodel.SlotsMeta{
		SlotID:           dto.SlotID,
		WorkerID:         dto.WorkerID,
		//WorkerScheduleID: dto.WorkerScheduleID,
	}
}

func SlotUpdateToModel(dto *orgdto.SlotUpdate) *orgmodel.SlotsMeta {
	return &orgmodel.SlotsMeta{
		SlotID:           dto.SlotID,
		WorkerID:         dto.WorkerID,
		//WorkerScheduleID: dto.WorkerScheduleID,
	}
}

func SlotInfoToDTO(model *orgmodel.Slot) *orgdto.Slot {
	return &orgdto.Slot{
		WorkerScheduleID: model.WorkerScheduleID,
		WorkerID:         model.WorkerID,
		Date:             strings.Fields(model.Date.String())[0],
		Begin:            model.Begin.Format(timeFormat),
		End:              model.End.Format(timeFormat),
		Busy:             model.Busy,
	}
}
