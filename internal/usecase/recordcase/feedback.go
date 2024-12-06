package recordcase

import (
	"context"
	"errors"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mapper/recordmap"

	"go.uber.org/zap"
)

func (r *RecordUseCase) Feedback(ctx context.Context, params *recordto.FeedbackParams) (*recordto.Feedback, error) {
	data, err := r.records.Feedback(ctx, recordmap.FeedParamsToModel(params))
	if err != nil {
		if errors.Is(err, postgres.ErrFeedbackNotFound) {
			return nil, err
		}
		r.Logger.Error(
			"failed to get feedback",
			zap.Error(err),
		)
		return nil, err
	}
	return recordmap.FeedbackToDTO(data), nil
}

func (r *RecordUseCase) FeedbackSet(ctx context.Context, feedback *recordto.Feedback) error {
	if err := r.records.FeedbackSet(ctx, recordmap.FeedbackToModel(feedback)); err != nil {
		r.Logger.Error(
			"failed to set feedback",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (r *RecordUseCase) FeedbackUpdate(ctx context.Context, feedback *recordto.Feedback) error {
	if err := r.records.FeedbackUpdate(ctx, recordmap.FeedbackToModel(feedback)); err != nil {
		r.Logger.Error(
			"failed to update feedback",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (r *RecordUseCase) FeedbackDelete(ctx context.Context, params *recordto.FeedbackParams) error {
	if err := r.records.FeedbackDelete(ctx, recordmap.FeedParamsToModel(params)); err != nil {
		r.Logger.Error(
			"failed to delete feedback",
			zap.Error(err),
		)
		return err
	}
	//
	return nil
}
