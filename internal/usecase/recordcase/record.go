package recordcase

import (
	"context"
	"errors"
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
}

func New(backdata *loader.BackData, userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, recordRepo infrastructure.RecordRepository, mailRepo infrastructure.Mail) *RecordUseCase {
	return &RecordUseCase{
		backdata: backdata,
		users:    userRepo,
		orgs:     orgRepo,
		records:  recordRepo,
		mail:     mailRepo,
	}
}

func (r *RecordUseCase) Record(ctx context.Context, logger *zap.Logger, recordID int) (*recordto.RecordScrap, error) {
	data, err := r.records.Record(ctx, recordID)
	if err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched record")
	return recordmap.RecordScrapToDTO(data), nil
}

func (r *RecordUseCase) RecordList(ctx context.Context, logger *zap.Logger, params *recordto.RecordListParams) (*recordto.RecordList, error) {
	data, found, err := r.records.RecordList(ctx, recordmap.RecordParamsToModel(params))
	if err != nil {
		if errors.Is(err, postgres.ErrRecordsNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched record list")
	resp := &recordto.RecordList{
		List:  recordmap.RecordListToDTO(data),
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
	logger.Info("Record has been saved")
	r.mail.SendMsg(&models.Message{
		Email:    record.UserEmail,
		Type:     mail.ReminderType,
		Value:    recordmap.RecordToReminder(record),
		IsAttach: true,
	})
	logger.Info("Notification has been sent to user's email")
	return nil
}

func (r *RecordUseCase) RecordPatch(ctx context.Context, logger *zap.Logger, rec *recordto.Record) error {
	if err := r.records.RecordPatch(ctx, recordmap.RecordToModel(rec)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Record has been patched")
	return nil
}

func (r *RecordUseCase) RecordCancel(ctx context.Context, logger *zap.Logger, rec *recordto.RecordCancelation) error {
	if err := r.records.RecordCancel(ctx, recordmap.CancelationToModel(rec)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Record has been canceled")
	return nil
}
