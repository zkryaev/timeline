package orgcase

import (
	"context"
	"errors"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/usecase/common"
	"timeline/internal/usecase/common/validation"

	"go.uber.org/zap"
)

func (o *OrgUseCase) TimetableAdd(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error {
	logger.Info("Checking timetable...")
	for i := range newTimetable.Timetable {
		if !validation.IsPeriodValid(newTimetable.Timetable[i].Open, newTimetable.Timetable[i].Close) {
			return common.ErrTimeIncorrect
		}
	}
	logger.Info("Timetable is valid")
	if err := o.org.TimetableAdd(ctx, newTimetable.OrgID, orgmap.TimetableToModel(newTimetable.Timetable)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Timetable has been saved")
	return nil
}

func (o *OrgUseCase) TimetableUpdate(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error {
	logger.Info("Checking timetable...")
	for i := range newTimetable.Timetable {
		if !validation.IsPeriodValid(newTimetable.Timetable[i].Open, newTimetable.Timetable[i].Close) {
			return common.ErrTimeIncorrect
		}
	}
	logger.Info("Timetable is valid")
	if err := o.org.TimetableUpdate(ctx, newTimetable.OrgID, orgmap.TimetableToModel(newTimetable.Timetable)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Timetable has been updated")
	return nil
}

func (o *OrgUseCase) TimetableDelete(ctx context.Context, logger *zap.Logger, orgID, weekday int) error {
	if err := o.org.TimetableDelete(ctx, orgID, weekday); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Timetable has been deleted")
	return nil
}

func (o *OrgUseCase) Timetable(ctx context.Context, logger *zap.Logger, req orgdto.TimetableReq) (*orgdto.Timetable, error) {
	timetable, city, err := o.org.Timetable(ctx, req.OrgID, req.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrTimetableNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched timetable")
	resp := &orgdto.Timetable{
		OrgID:     req.OrgID,
		Timetable: make([]*entity.OpenHours, 0, len(timetable)),
	}
	tzid := o.backdata.Cities.GetCityTZ(city)
	loc, err := time.LoadLocation(tzid)
	if err != nil {
		logger.Error("failed to load location, set UTC+03 (MSK)", zap.String("city-tzid", city+"="+tzid), zap.Error(err))
		loc = time.Local // UTC+03 = MSK
	}
	for _, v := range timetable {
		resp.Timetable = append(resp.Timetable, orgmap.OpenHoursToDTO(v, loc))
	}
	return resp, nil
}
