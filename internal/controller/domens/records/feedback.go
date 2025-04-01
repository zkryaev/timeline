package records

import (
	"context"
	"net/http"
	"strconv"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/libs/custom"

	"go.uber.org/zap"
)

type Feedback interface {
	FeedbackList(ctx context.Context, params *recordto.FeedbackParams) (*recordto.FeedbackList, error)
	FeedbackSet(ctx context.Context, feedback *recordto.Feedback) error
	FeedbackUpdate(ctx context.Context, feedback *recordto.Feedback) error
	FeedbackDelete(ctx context.Context, params *recordto.FeedbackParams) error
}

// @Summary Feedback info
// @Description Get feedbakc for specified record
// @Tags record / feedback
// @Param limit query int true "Limit the number of results"
// @Param page query int true "Page number for pagination"
// @Param   recordID query int true "record_id"
// @Param   orgID query int true "org_id"
// @Param   userID query int true "user_id"
// @Success 200 {object} recordto.FeedbackList
// @Failure 400
// @Failure 500
// @Router /records/feedbacks/info [get]
func (rec *RecordCtrl) Feedbacks(w http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	query := map[string]bool{
		"limit":     true,
		"page":      true,
		"record_id": false,
		"org_id":    false,
		"user_id":   false,
	}
	if err := validation.IsQueryValid(r, query); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	params := map[string]string{
		"limit":     "int",
		"page":      "int",
		"record_id": "int",
		"org_id":    "int",
		"user_id":   "int",
	}
	queryParams, err := custom.QueryParamsConv(params, r.URL.Query())
	if err != nil {
		logger.Error("QueryParamsConv", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &recordto.FeedbackParams{
		RecordID: queryParams["record_id"].(int),
		UserID:   queryParams["user_id"].(int),
		OrgID:    queryParams["org_id"].(int),
		Limit:    queryParams["limit"].(int),
		Page:     queryParams["page"].(int),
	}
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := rec.usecase.FeedbackList(r.Context(), req)
	if err != nil {
		logger.Error("FeedbackList", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := common.WriteJSON(w, data); err != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Set feedback
// @Description Set feedback for specified record
// @Tags record / feedback
// @Accept  json
// @Param feedback body recordto.Feedback true "Feedback data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [post]
func (rec *RecordCtrl) FeedbackSet(w http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	req := &recordto.Feedback{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackSet(r.Context(), req); err != nil {
		logger.Error("FeedbackSet", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update feedback
// @Description Update feedback for specified record
// @Tags record / feedback
// @Accept  json
// @Param feedback body recordto.Feedback true "Feedback data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [put]
func (rec *RecordCtrl) FeedbackUpdate(w http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	req := &recordto.Feedback{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackUpdate(r.Context(), req); err != nil {
		logger.Error("FeedbackUpdate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete feedback
// @Description Delete feedback for specified record
// @Tags record / feedback
// @Param   recordID query int true "record_id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/feedbacks/info [delete]
func (rec *RecordCtrl) FeedbackDelete(w http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	query := map[string]bool{
		"record_id": true,
	}
	if err := validation.IsQueryValid(r, query); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	recordid, err := strconv.Atoi(r.URL.Query().Get("record_id"))
	if err != nil {
		logger.Error("Atoi", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &recordto.FeedbackParams{RecordID: recordid}
	if common.Validate(req) != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err = rec.usecase.FeedbackDelete(r.Context(), req); err != nil {
		logger.Error("FeedbackDelete", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
