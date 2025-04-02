package orgcase

import (
	"context"
	"errors"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/usecase/common"
	"timeline/internal/usecase/common/validation"

	"go.uber.org/zap"
)

func (o *OrgUseCase) WorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) (*orgdto.ScheduleList, error) {
	data, err := o.org.WorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params))
	if err != nil {
		if errors.Is(err, postgres.ErrScheduleNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched worker schedule")
	return orgmap.ScheduleListToDTO(data), nil
}

func (o *OrgUseCase) AddWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error {
	logger.Info("Checking worker schedule...")
	for i := range schedule.Schedule {
		if !validation.IsPeriodValid(schedule.Schedule[i].Start, schedule.Schedule[i].Over) {
			return common.ErrTimeIncorrect
		}
	}
	logger.Info("Worker schedule is valid")
	if err := o.org.AddWorkerSchedule(ctx, orgmap.WorkerScheduleToModel(schedule)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Worker schedule has been saved")
	return nil
}

func (o *OrgUseCase) UpdateWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error {
	logger.Info("Worker session duration has been updated")
	if err := o.org.UpdateWorkerSchedule(ctx, orgmap.WorkerScheduleToModel(schedule)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Worker schedule has been updated")
	return nil
}

func (o *OrgUseCase) DeleteWorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) error {
	if err := o.org.SoftDeleteWorkerSchedule(ctx, orgmap.ScheduleParamsToModel(params)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Worker schedule has been deleted")
	return nil
}
