package postgres

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/mapper/recordmap"
)

func (suite *PostgresTestSuite) TestFeedbackQueries() {
	ctx := context.Background()

	recordID := 2
	recScrap, err := suite.db.Record(ctx, recordID)
	suite.Require().NoError(err)

	exp := &recordto.Feedback{
		RecordID:   recordID,
		Stars:      4,
		Feedback:   "Хорошая тренировка, но хотелось бы больше внимания.",
		Service:    recScrap.Service.Name,
		FirstName:  recScrap.Worker.FirstName,
		LastName:   recScrap.Worker.LastName,
		RecordDate: recScrap.CreatedAt.Format(time.DateOnly),
	}

	suite.Require().NoError(suite.db.FeedbackSet(ctx, recordmap.FeedbackToModel(exp)), fmt.Sprintf("record_id=%d", exp.RecordID))

	params := &recordto.FeedbackParams{
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
			suite.Equal(exp, feedback)
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
			suite.Equal(exp, feedback)
		}
	}
	params.RecordID = exp.RecordID
	suite.NoError(suite.db.FeedbackDelete(ctx, recordmap.FeedParamsToModel(params)))
}
