package orgcase

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mapper/orgmap"

	"go.uber.org/zap"
)

func (o *OrgUseCase) TimetableAdd(ctx context.Context, newTimetable *orgdto.Timetable) error {
	// Валидация начала и конца работы организации
	for i := range newTimetable.Timetable {
		if !workPeriodValid(newTimetable.Timetable[i].Open, newTimetable.Timetable[i].Close) {
			o.Logger.Error(
				"failed cause given time is somehow incorrect",
			)
			return fmt.Errorf("some of the provided time is incorrect")
		}
	}
	if err := o.org.TimetableAdd(ctx, newTimetable.OrgID, orgmap.TimetableToModel(newTimetable.Timetable)); err != nil {
		if errors.Is(err, postgres.ErrOrgNotFound) {
			return err
		}
		o.Logger.Error(
			"failed to add org timetable",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) TimetableUpdate(ctx context.Context, newTimetable *orgdto.Timetable) error {
	for i := range newTimetable.Timetable {
		if !workPeriodValid(newTimetable.Timetable[i].Open, newTimetable.Timetable[i].Close) {
			o.Logger.Error(
				"failed cause given time is somehow incorrect",
			)
			return fmt.Errorf("some of the provided time is incorrect")
		}
	}
	if err := o.org.TimetableUpdate(ctx, newTimetable.OrgID, orgmap.TimetableToModel(newTimetable.Timetable)); err != nil {
		if errors.Is(err, postgres.ErrOrgNotFound) {
			return err
		}
		o.Logger.Error(
			"failed to update org timetable",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (o *OrgUseCase) TimetableDelete(ctx context.Context, orgID, weekday int) error {
	if err := o.org.TimetableDelete(ctx, orgID, weekday); err != nil {
		if errors.Is(err, postgres.ErrOrgNotFound) {
			return err
		}
		o.Logger.Error(
			"failed to delete org timetable",
			zap.Error(err),
		)
		return err
	}
	return nil
}
