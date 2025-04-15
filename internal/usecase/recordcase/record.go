package recordcase

import (
	"context"
	"errors"
	"time"
	"timeline/internal/controller/scope"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mail"
	"timeline/internal/infrastructure/mapper/recordmap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/usecase/common"
	"timeline/internal/utils/loader"

	"go.uber.org/zap"
)

type RecordUseCase struct {
	backdata *loader.BackData
	users    infrastructure.UserRepository
	orgs     infrastructure.OrgRepository
	records  infrastructure.RecordRepository
	mail     infrastructure.Mail
	settings *scope.Settings
}

func New(backdata *loader.BackData, userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository,
	recordRepo infrastructure.RecordRepository, mailRepo infrastructure.Mail, settings *scope.Settings) *RecordUseCase {
	return &RecordUseCase{
		backdata: backdata,
		users:    userRepo,
		orgs:     orgRepo,
		records:  recordRepo,
		mail:     mailRepo,
		settings: settings,
	}
}

func (r *RecordUseCase) Record(ctx context.Context, logger *zap.Logger, param recordto.RecordParam) (*recordto.RecordScrap, error) {
	data, err := r.records.Record(ctx, recordmap.RecordParamToModel(param))
	if err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	var city string
	if param.TData.IsOrg {
		city = data.User.City
	} else {
		city = data.Org.City
	}
	tzid := r.backdata.Cities.GetCityTZ(city)
	loc, err := time.LoadLocation(tzid)
	if err != nil {
		logger.Error("failed to load location, set UTC+03 (MSK)", zap.String("city-tzid", data.User.City+"="+tzid), zap.Error(err))
		loc = time.Local // UTC+03 = MSK
	}
	logger.Info("Fetched record")
	return recordmap.RecordScrapToDTO(data, loc), nil
}

func (r *RecordUseCase) RecordList(ctx context.Context, logger *zap.Logger, params *recordto.RecordListParams) (*recordto.RecordList, error) {
	if !params.TData.IsOrg {
		params.UserID = params.TData.ID
	} else {
		params.OrgID = params.TData.ID
	}
	data, found, err := r.records.RecordList(ctx, recordmap.RecordParamsToModel(params))
	if err != nil {
		if errors.Is(err, postgres.ErrRecordsNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched record list")
	var loc *time.Location
	if len(data) > 0 {
		var city string
		if params.TData.IsOrg {
			for i := range data {
				if data[i].Org.OrgID == params.TData.ID {
					city = data[i].Org.City
					break
				}
			}
		} else {
			city = data[0].User.City
		}
		tzid := r.backdata.Cities.GetCityTZ(city)
		loc, err = time.LoadLocation(tzid)
		if err != nil {
			logger.Error("failed to load location, set UTC+03 (MSK)", zap.String("city-tzid", data[0].User.City+"="+tzid), zap.Error(err))
			loc = time.Local
		}
	} else {
		return &recordto.RecordList{Found: found}, nil
	}
	resp := &recordto.RecordList{
		List:  recordmap.RecordListToDTO(data, loc),
		Found: found,
	}
	return resp, nil
}

func (r *RecordUseCase) RecordAdd(ctx context.Context, logger *zap.Logger, rec *recordto.Record) error {
	record, _, err := r.records.RecordAdd(ctx, recordmap.RecordToModel(rec))
	if err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	tzid := r.backdata.Cities.GetCityTZ(record.UserCity)
	loc, err := time.LoadLocation(tzid)
	if err != nil {
		logger.Error("failed to load location, set UTC+03 (MSK)", zap.String("city-tzid", record.UserCity+"="+tzid), zap.Error(err))
		loc = time.Local
	}

	logger.Info("Record has been saved")
	if r.settings.EnableRepoS3 {
		r.mail.SendMsg(&models.Message{
			Email:    record.UserEmail,
			Type:     mail.CancelationType,
			Value:    recordmap.ReminderRecordToReminder(record, loc),
			IsAttach: true,
		})
		logger.Info("Notification has been sent to user's email")
	}
	return nil
}

func (r *RecordUseCase) RecordCancel(ctx context.Context, logger *zap.Logger, req *recordto.RecordCancelation) error {
	param := recordto.RecordParam{RecordID: req.RecordID, TData: req.TData}
	record, err := r.Record(ctx, logger, param)
	if err != nil {
		return err
	}
	if err := r.records.RecordCancel(ctx, recordmap.CancelationToModel(req)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Record has been canceled")
	if r.settings.EnableRepoS3 {
		r.mail.SendMsg(&models.Message{
			Email: record.User.Email,
			Type:  mail.ReminderType,
			Value: models.CancelMsg{
				Organization: record.Org.Info.Name,
				Service:      record.Service.Name,
				ServiceDecs:  record.Service.Description,
				SessionStart: record.Slot.Begin,
				SessionEnd:   record.Slot.End,
				SessionDate:  record.Slot.Date,
				CancelReason: req.CancelReason,
			},
		})
		logger.Info("Notification has been sent to user's email")
	}
	return nil
}
