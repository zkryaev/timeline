package records

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/query"
	"timeline/internal/controller/scope"
	"timeline/internal/entity/dto/recordto"

	"go.uber.org/zap"
)

type Feedback interface {
	FeedbackList(ctx context.Context, logger *zap.Logger, params *recordto.FeedbackParams) (*recordto.FeedbackList, error)
	FeedbackSet(ctx context.Context, logger *zap.Logger, feedback *recordto.Feedback) error
	FeedbackUpdate(ctx context.Context, logger *zap.Logger, feedback *recordto.Feedback) error
	FeedbackDelete(ctx context.Context, logger *zap.Logger, params *recordto.FeedbackParams) error
}

// @Summary Feedback info
// @Description Get feedback for specified record
// @Tags records/feedbacks
// @Param limit query int true " "
// @Param page query int true " "
// @Param record_id query int false " "
// @Param org_id query int false " "
// @Param user_id query int false " "
// @Success 200 {object} recordto.FeedbackList
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [get]
func (rec *RecordCtrl) Feedbacks(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(rec.settings, rec.Logger, r.Context())
	var (
		limit    = query.NewParamInt(scope.LIMIT, true)
		page     = query.NewParamInt(scope.PAGE, true)
		orgID    = query.NewParamInt(scope.ORG_ID, false)
		userID   = query.NewParamInt(scope.USER_ID, false)
		recordID = query.NewParamInt(scope.RECORD_ID, false)
	)
	params := query.NewParams(rec.settings, limit, page, orgID, userID, recordID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &recordto.FeedbackParams{
		RecordID: recordID.Val,
		UserID:   userID.Val,
		OrgID:    orgID.Val,
		Limit:    limit.Val,
		Page:     page.Val,
	}
	data, err := rec.usecase.FeedbackList(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("FeedbackList", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("FeedbackList", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	if err := common.WriteJSON(w, data); err != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Set feedback
// @Description Set feedback for specified record
// @Tags records/feedbacks
// @Accept  json
// @Param req body recordto.Feedback true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [post]
func (rec *RecordCtrl) FeedbackSet(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(rec.settings, rec.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(rec.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &recordto.Feedback{TData: tdata}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackSet(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("FeedbackSet", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("FeedbackSet", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update feedback
// @Description Update feedback for specified record
// @Tags records/feedbacks
// @Accept  json
// @Param req body recordto.Feedback true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [put]
func (rec *RecordCtrl) FeedbackUpdate(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(rec.settings, rec.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(rec.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &recordto.Feedback{TData: tdata}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackUpdate(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("FeedbackUpdate", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("FeedbackUpdate", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete feedback
// @Description Delete feedback for specified record
// @Tags records/feedbacks
// @Param   record_id query int true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [delete]
func (rec *RecordCtrl) FeedbackDelete(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(rec.settings, rec.Logger, r.Context())
	var (
		recordID = query.NewParamInt(scope.RECORD_ID, true)
	)
	params := query.NewParams(rec.settings, recordID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &recordto.FeedbackParams{RecordID: recordID.Val}
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackDelete(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("FeedbackDelete", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("FeedbackDelete", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
