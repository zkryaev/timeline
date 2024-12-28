package recordcase

import (
	"context"
	"errors"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mapper/recordmap"

	"go.uber.org/zap"
)

func (r *RecordUseCase) FeedbackList(ctx context.Context, params *recordto.FeedbackParams) (*recordto.FeedbackList, error) {
	data, found, err := r.records.FeedbackList(ctx, recordmap.FeedParamsToModel(params))
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
	feedbackList := make([]*recordto.Feedback, 0, 1)
	for i := range data {
		feedbackList = append(feedbackList, recordmap.FeedbackToDTO(data[i]))
	}
	resp := &recordto.FeedbackList{
		List:  feedbackList,
		Found: found,
	}
	return resp, nil
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
