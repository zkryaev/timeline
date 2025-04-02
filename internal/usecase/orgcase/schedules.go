package orgcase

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

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

func (o *OrgUseCase) WorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) (*orgdto.ScheduleList, error) {
	data, err := o.org.WorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params))
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched worker schedule")
	return orgmap.ScheduleListToDTO(data), nil
}

func (o *OrgUseCase) AddWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error {
	logger.Info("Checking worker schedule...")
	for i := range schedule.Schedule {
		if !workPeriodValid(schedule.Schedule[i].Start, schedule.Schedule[i].Over) {
			return fmt.Errorf("some of the provided time is incorrect")
		}
	}
	logger.Info("Worker schedule is valid")
	if err := o.org.AddWorkerSchedule(ctx, orgmap.WorkerScheduleToModel(schedule)); err != nil {
		return err
	}
	logger.Info("Worker schedule has been saved")
	return nil
}

func (o *OrgUseCase) UpdateWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error {
	worker := &orgdto.UpdateWorkerReq{
		WorkerID: schedule.WorkerID,
		OrgID:    schedule.OrgID,
		WorkerInfo: entity.Worker{
			SessionDuration: schedule.SessionDuration,
		},
	}
	if err := o.WorkerPatch(ctx, logger, worker); err != nil {
		return err
	}
	logger.Info("Worker session duration has been updated")
	if err := o.org.UpdateWorkerSchedule(ctx, orgmap.WorkerScheduleToModel(schedule)); err != nil {
		return err
	}
	logger.Info("Worker schedule has been updated")
	return nil
}

func (o *OrgUseCase) DeleteWorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) error {
	if err := o.org.SoftDeleteWorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params)); err != nil {
		return err
	}
	logger.Info("Worker schedule has been deleted")
	return nil
}
