package recordcase

import (
	"context"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/mail"
	"timeline/internal/infrastructure/mapper/recordmap"
	"timeline/internal/infrastructure/models"

	"go.uber.org/zap"
)

type RecordUseCase struct {
	users   infrastructure.UserRepository
	orgs    infrastructure.OrgRepository
	records infrastructure.RecordRepository
	mail    infrastructure.Mail
	Logger  *zap.Logger
}

func New(userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, recordRepo infrastructure.RecordRepository, mailRepo infrastructure.Mail, logger *zap.Logger) *RecordUseCase {
	return &RecordUseCase{
		users:   userRepo,
		orgs:    orgRepo,
		records: recordRepo,
		mail:    mailRepo,
		Logger:  logger,
	}
}

func (r *RecordUseCase) Record(ctx context.Context, recordID int) (*recordto.RecordScrap, error) {
	data, err := r.records.Record(ctx, recordID)
	if err != nil {
		r.Logger.Error(
			"failed to get record",
			zap.Error(err),
		)
		return nil, err
	}
	return recordmap.RecordScrapToDTO(data), nil
}

func (r *RecordUseCase) RecordList(ctx context.Context, params *recordto.RecordListParams) (*recordto.RecordList, error) {
	data, found, err := r.records.RecordList(ctx, recordmap.RecordParamsToModel(params))
	if err != nil {
		r.Logger.Error(
			"failed to get record list",
			zap.Error(err),
		)
		return nil, err
	}
	resp := &recordto.RecordList{
		List:  recordmap.RecordListToDTO(data),
		Found: found,
	}
	return resp, nil
}

func (r *RecordUseCase) RecordAdd(ctx context.Context, rec *recordto.Record) error {
	record, _, err := r.records.RecordAdd(ctx, recordmap.RecordToModel(rec))
	if err != nil {
		r.Logger.Error(
			"failed to add record",
			zap.Error(err),
		)
		return err
	}

	r.mail.SendMsg(&models.Message{
		Email:    record.UserEmail,
		Type:     mail.ReminderType,
		Value:    recordmap.RecordToReminder(record),
		IsAttach: true,
	})
	return nil
}

func (r *RecordUseCase) RecordPatch(ctx context.Context, rec *recordto.Record) error {
	if err := r.records.RecordPatch(ctx, recordmap.RecordToModel(rec)); err != nil {
		r.Logger.Error(
			"failed to add record",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (r *RecordUseCase) RecordDelete(ctx context.Context, rec *recordto.Record) error {
	if err := r.records.RecordSoftDelete(ctx, rec.RecordID); err != nil {
		r.Logger.Error(
			"failed to add record",
			zap.Error(err),
		)
		return err
	}
	return nil
}
