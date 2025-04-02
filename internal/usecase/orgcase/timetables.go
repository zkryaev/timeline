package orgcase

import (
	"context"
	"fmt"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) TimetableAdd(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error {
	logger.Info("Checking timetable...")
	for i := range newTimetable.Timetable {
		if !workPeriodValid(newTimetable.Timetable[i].Open, newTimetable.Timetable[i].Close) {
			return fmt.Errorf("some of the provided time is incorrect")
		}
	}
	logger.Info("Timetable is valid")
	if err := o.org.TimetableAdd(ctx, newTimetable.OrgID, orgmap.TimetableToModel(newTimetable.Timetable)); err != nil {
		return err
	}
	logger.Info("Timetable has been saved")
	return nil
}

func (o *OrgUseCase) TimetableUpdate(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error {
	logger.Info("Checking timetable...")
	for i := range newTimetable.Timetable {
		if !workPeriodValid(newTimetable.Timetable[i].Open, newTimetable.Timetable[i].Close) {
			return fmt.Errorf("some of the provided time is incorrect")
		}
	}
	logger.Info("Timetable is valid")
	if err := o.org.TimetableUpdate(ctx, newTimetable.OrgID, orgmap.TimetableToModel(newTimetable.Timetable)); err != nil {
		return err
	}
	logger.Info("Timetable has been updated")
	return nil
}

func (o *OrgUseCase) TimetableDelete(ctx context.Context, logger *zap.Logger, orgID, weekday int) error {
	if err := o.org.TimetableDelete(ctx, orgID, weekday); err != nil {
		return err
	}
	logger.Info("Timetable has been deleted")
	return nil
}

func (o *OrgUseCase) Timetable(ctx context.Context, logger *zap.Logger, orgID int) (*orgdto.Timetable, error) {
	timetable, err := o.org.Timetable(ctx, orgID)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched timetable")
	resp := &orgdto.Timetable{
		OrgID:     orgID,
		Timetable: make([]*entity.OpenHours, 0, len(timetable)),
	}
	for _, v := range timetable {
		resp.Timetable = append(resp.Timetable, orgmap.OpenHoursToDTO(v))
	}
	return resp, nil
}
