package postgres

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/mapper/recordmap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/recordmodel"
)

func (suite *PostgresTestSuite) TestFeedbackQueries() {
	ctx := context.Background()
	tdata := models.TokenData{ID: 1, IsOrg: false}
	recordreq := recordmodel.RecordParam{RecordID: 2, TData: tdata}
	recScrap, err := suite.db.Record(ctx, recordreq)
	suite.Require().NoError(err)

	exp := &recordto.Feedback{
		TData:           entity.TokenData(tdata),
		RecordID:        recordreq.RecordID,
		Stars:           4,
		Feedback:        "Хорошая тренировка, но хотелось бы больше внимания.",
		Service:         recScrap.Service.Name,
		WorkerFirstName: recScrap.Worker.FirstName,
		WorkerLastName:  recScrap.Worker.LastName,
		UserFirstName:   recScrap.User.FirstName,
		UserLastName:    recScrap.User.LastName,
		RecordDate:      recScrap.CreatedAt.Format(time.DateOnly),
	}

	suite.Require().NoError(suite.db.FeedbackSet(ctx, recordmap.FeedbackToModel(exp)), fmt.Sprintf("record_id=%d", exp.RecordID))

	params := &recordto.FeedbackParams{
		TData: entity.TokenData(tdata),
		Limit: 5,
		Page:  1,
	}

	feedbkList, found, err := suite.db.FeedbackList(ctx, recordmap.FeedParamsToModel(params))
	suite.NoError(err)
	suite.NotZero(found)
	suite.NotNil(feedbkList)

	for _, feedbk := range feedbkList {
		feedback := recordmap.FeedbackToDTO(feedbk)
		if exp.RecordID == feedback.RecordID {
			suite.Equal(exp.Feedback, feedback.Feedback)
			suite.Equal(exp.Stars, feedback.Stars)
		}
	}
	exp.Feedback = "ТЕСТИРОВАНИЕ"
	suite.NoError(suite.db.FeedbackUpdate(ctx, recordmap.FeedbackToModel(exp)))

	feedbkList, found, err = suite.db.FeedbackList(ctx, recordmap.FeedParamsToModel(params))
	suite.NoError(err)
	suite.NotZero(found)
	suite.NotNil(feedbkList)

	for _, feedbk := range feedbkList {
		feedback := recordmap.FeedbackToDTO(feedbk)
		if exp.RecordID == feedback.RecordID {
			suite.Equal(exp.Feedback, feedback.Feedback)
		}
	}
	params.RecordID = exp.RecordID
	suite.NoError(suite.db.FeedbackDelete(ctx, recordmap.FeedParamsToModel(params)))
}
