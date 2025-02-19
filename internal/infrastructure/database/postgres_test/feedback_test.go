package postgres_test

import (
	"context"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/infrastructure/mapper/recordmap"
)

func (suite *PostgresTestSuite) TestFeedbackQueries() {
	ctx := context.Background()

	// (2, 4, 'Хорошая тренировка, но хотелось бы больше внимания.');
	exp := &recordto.Feedback{
		RecordID: 2, // хардкод!
		Stars:    4,
		Feedback: "Хорошая тренировка, но хотелось бы больше внимания.",
	}

	suite.NoError(suite.db.FeedbackSet(ctx, recordmap.FeedbackToModel(exp)))

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
