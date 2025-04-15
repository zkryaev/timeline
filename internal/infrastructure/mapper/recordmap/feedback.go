package recordmap

import (
	"database/sql"
	"time"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/recordmodel"
)

func FeedbackToModel(dto *recordto.Feedback) *recordmodel.Feedback {
	return &recordmodel.Feedback{
		FeedbackID: sql.NullInt32{Int32: int32(dto.FeedbackID), Valid: true},
		RecordID:   sql.NullInt32{Int32: int32(dto.RecordID), Valid: true},
		Stars:      sql.NullInt32{Int32: int32(dto.Stars), Valid: true},
		Feedback:   sql.NullString{String: dto.Feedback, Valid: true},
		TData:      models.TokenData(dto.TData),
	}
}

func FeedbackToDTO(model *recordmodel.Feedback) *recordto.Feedback {
	if model.Stars.Valid {
		return &recordto.Feedback{
			FeedbackID:      int(model.FeedbackID.Int32),
			RecordID:        int(model.RecordID.Int32),
			Stars:           int(model.Stars.Int32),
			Feedback:        model.Feedback.String,
			Service:         model.Service.String,
			WorkerFirstName: model.WorkerFirstName.String,
			WorkerLastName:  model.WorkerLastName.String,
			UserFirstName:   model.UserFirstName.String,
			UserLastName:    model.UserLastName.String,
			RecordDate:      model.RecordDate.Time.Format(time.DateOnly),
		}
	}
	return nil
}

func FeedParamsToModel(dto *recordto.FeedbackParams) *recordmodel.FeedbackParams {
	return &recordmodel.FeedbackParams{
		FeedbackID: dto.FeedbackID,
		RecordID:   dto.RecordID,
		UserID:     dto.UserID,
		OrgID:      dto.OrgID,
		Limit:      dto.Limit,
		Offset:     (dto.Page - 1) * dto.Limit,
		TData:      models.TokenData(dto.TData),
	}
}
