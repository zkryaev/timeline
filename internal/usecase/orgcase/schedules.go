package orgcase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mapper/orgmap"

	"go.uber.org/zap"
)

// Если over >= start - false
func workPeriodValid(start, over string) bool {
	overTime, errover := time.Parse("15:06", over)
	startTime, errstart := time.Parse("15:06", start)
	if overTime.Compare(startTime) <= 0 || errover != nil || errstart != nil {
		return false
	}
	return true
}

func (o *OrgUseCase) WorkerSchedule(ctx context.Context, params *orgdto.ScheduleParams) (*orgdto.ScheduleList, error) {
	data, err := o.org.WorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params))
	if err != nil {
		o.Logger.Error(
			"failed to get worker schedule",
			zap.Error(err),
		)
		return nil, err
	}
	return orgmap.ScheduleListToDTO(data), nil
}

func (o *OrgUseCase) AddWorkerSchedule(ctx context.Context, schedule *orgdto.ScheduleList) error {
	for i := range schedule.Schedule {
		if !workPeriodValid(schedule.Schedule[i].Start, schedule.Schedule[i].Over) {
			o.Logger.Error(
				"failed cause given time is somehow incorrect",
			)
			return fmt.Errorf("some of the provided time is incorrect")
		}
	}
	if err := o.org.AddWorkerSchedule(ctx, orgmap.ScheduleListToModel(schedule)); err != nil {
		o.Logger.Error(
			"failed to get worker schedule",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) UpdateWorkerSchedule(ctx context.Context, schedule *orgdto.ScheduleList) error {
	worker := &orgdto.UpdateWorkerReq{
		WorkerID: schedule.WorkerID,
		OrgID:    schedule.OrgID,
		WorkerInfo: entity.Worker{
			SessionDuration: schedule.SessionDuration,
		},
	}
	if err := o.WorkerPatch(ctx, worker); err != nil {
		return err
	}
	if err := o.org.UpdateWorkerSchedule(ctx, orgmap.ScheduleListToModel(schedule)); err != nil {
		o.Logger.Error(
			"failed to get worker schedule",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) DeleteWorkerSchedule(ctx context.Context, params *orgdto.ScheduleParams) error {
	if err := o.org.DeleteWorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params)); err != nil {
		if errors.Is(err, postgres.ErrScheduleNotFound) {
			return err
		}
		o.Logger.Error(
			"failed to get worker schedule",
			zap.Error(err),
		)
		return err
	}
	return nil
}
