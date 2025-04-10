package records

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/sugar/custom"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Record interface {
	Record(ctx context.Context, logger *zap.Logger, recordID int) (*recordto.RecordScrap, error)
	RecordList(ctx context.Context, logger *zap.Logger, params *recordto.RecordListParams) (*recordto.RecordList, error)
	RecordAdd(ctx context.Context, logger *zap.Logger, rec *recordto.Record) error
	RecordCancel(ctx context.Context, logger *zap.Logger, rec *recordto.RecordCancelation) error
	Feedback
}

type RecordCtrl struct {
	usecase    Record
	Logger     *zap.Logger
	middleware middleware.Middleware
}

func New(usecase Record, middleware middleware.Middleware, logger *zap.Logger) *RecordCtrl {
	return &RecordCtrl{
		usecase:    usecase,
		Logger:     logger,
		middleware: middleware,
	}
}

// @Summary Record information
// @Description Get bounded with record information
// @Tags Records
// @Param   recordID path int true "record_id"
// @Success 200 {object} recordto.RecordScrap
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /records/info/{recordID} [get]
func (rec *RecordCtrl) Record(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "recordID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if params["recordID"] <= 0 {
		logger.Error("record_id must be > 0", zap.Int("record_id", params["recordID"]))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := rec.usecase.Record(r.Context(), logger, params["recordID"])
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Record", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Record", zap.Error(err))
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

// @Summary Records
// @Description Get bounded with records informations
// @Tags Records
// @Param limit query int true "Limit the number of results"
// @Param page query int true "Page number for pagination"
// @Param   userID query int false "user_id"
// @Param   orgID query int false "org_id"
// @Param   fresh query bool false "Decide which records must be returned. True - only current & future records. False/NotGiven - olds"
// @Success 200 {object} recordto.RecordList
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /records/list [get]
func (rec *RecordCtrl) RecordList(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	query := map[string]bool{
		"user_id": false,
		"org_id":  false,
		"fresh":   false,
		"limit":   false,
		"page":    false,
	}
	if err := validation.IsQueryValid(r, query); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	params := map[string]string{
		"user_id": "int",
		"org_id":  "int",
		"fresh":   "bool",
		"limit":   "int",
		"page":    "int",
	}
	queryParams, err := custom.QueryParamsConv(params, r.URL.Query())
	if err != nil {
		logger.Error("QueryParamsConv", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &recordto.RecordListParams{
		UserID: queryParams["user_id"].(int),
		OrgID:  queryParams["org_id"].(int),
		Fresh:  queryParams["fresh"].(bool),
		Limit:  queryParams["limit"].(int),
		Page:   queryParams["page"].(int),
	}
	token, _ := rec.middleware.ExtractToken(r)
	tdata := common.GetTokenData(token.Claims)
	if !tdata.IsOrg {
		req.UserID = tdata.ID
	}
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := rec.usecase.RecordList(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("RecordList", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("RecordList", zap.Error(err))
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

// @Summary Add record
// @Description Add record with id components of other order details
// @Tags Records
// @Accept  json
// @Param record body recordto.Record true "Record data"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /records/creation [post]
func (rec *RecordCtrl) RecordAdd(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	req := &recordto.Record{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := rec.usecase.RecordAdd(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("RecordList", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("RecordList", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a future record
// @Description Cancel a future record. If record was done, it couldn't be deleted!
// @Tags Records
// @Param   cancelReq body recordto.RecordCancelation true "cancel description"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /records/info/{recordID} [put]
// Удаление только ожидаемой записи, а не уже совершённой.
func (rec *RecordCtrl) RecordCancel(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := rec.Logger.With(zap.String("uuid", uuid))
	req := &recordto.RecordCancelation{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := rec.usecase.RecordCancel(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("RecordCancel", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("RecordCancel", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
