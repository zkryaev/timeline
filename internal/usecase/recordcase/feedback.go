package recordcase

import (
	"context"
	"errors"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mapper/recordmap"
	"timeline/internal/usecase/common"

	"go.uber.org/zap"
)

func (r *RecordUseCase) FeedbackList(ctx context.Context, logger *zap.Logger, params *recordto.FeedbackParams) (*recordto.FeedbackList, error) {
	data, found, err := r.records.FeedbackList(ctx, recordmap.FeedParamsToModel(params))
	if err != nil {
		if errors.Is(err, postgres.ErrFeedbackNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched feedback list")
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

func (r *RecordUseCase) FeedbackSet(ctx context.Context, logger *zap.Logger, feedback *recordto.Feedback) error {
	if err := r.records.FeedbackSet(ctx, recordmap.FeedbackToModel(feedback)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Feedback has been saved")
	return nil
}

func (r *RecordUseCase) FeedbackUpdate(ctx context.Context, logger *zap.Logger, feedback *recordto.Feedback) error {
	if err := r.records.FeedbackUpdate(ctx, recordmap.FeedbackToModel(feedback)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Feedback has been updated")
	return nil
}

func (r *RecordUseCase) FeedbackDelete(ctx context.Context, logger *zap.Logger, params *recordto.FeedbackParams) error {
	if err := r.records.FeedbackDelete(ctx, recordmap.FeedParamsToModel(params)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Feedback has been deleted")
	return nil
}
