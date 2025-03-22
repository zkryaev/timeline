package records

import (
	"context"
	"net/http"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/libs/custom"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Record interface {
	Record(ctx context.Context, recordID int) (*recordto.RecordScrap, error)
	RecordList(ctx context.Context, params *recordto.RecordListParams) (*recordto.RecordList, error)
	RecordAdd(ctx context.Context, rec *recordto.Record) error
	RecordCancel(ctx context.Context, rec *recordto.RecordCancelation) error
	Feedback
}

type RecordCtrl struct {
	usecase Record
	Logger  *zap.Logger
}

func New(usecase Record, logger *zap.Logger) *RecordCtrl {
	return &RecordCtrl{
		usecase: usecase,
		Logger:  logger,
	}
}

// @Summary Record information
// @Description Get bounded with record information
// @Tags Records
// @Param   recordID path int true "record_id"
// @Success 200 {object} recordto.RecordScrap
// @Failure 400
// @Failure 500
// @Router /records/info/{recordID} [get]
func (rec *RecordCtrl) Record(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchPathID(mux.Vars(r), "recordID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if params["recordID"] <= 0 {
		http.Error(w, "record_id must be > 0", http.StatusBadRequest)
		return
	}
	data, err := rec.usecase.Record(r.Context(), params["recordID"])
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if common.WriteJSON(w, data) != nil {
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
// @Failure 500
// @Router /records/list [get]
func (rec *RecordCtrl) RecordList(w http.ResponseWriter, r *http.Request) {
	query := map[string]bool{
		"user_id": false,
		"org_id":  false,
		"fresh":   false,
		"limit":   false,
		"page":    false,
	}
	if !validation.IsQueryValid(r, query) {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
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
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}
	req := &recordto.RecordListParams{
		UserID: queryParams["user_id"].(int),
		OrgID:  queryParams["org_id"].(int),
		Fresh:  queryParams["fresh"].(bool),
		Limit:  queryParams["limit"].(int),
		Page:   queryParams["page"].(int),
	}
	if common.Validate(req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := rec.usecase.RecordList(r.Context(), req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if common.WriteJSON(w, data) != nil {
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
// @Failure 500
// @Router /records/creation [post]
func (rec *RecordCtrl) RecordAdd(w http.ResponseWriter, r *http.Request) {
	req := &recordto.Record{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := rec.usecase.RecordAdd(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a future record
// @Description Cancel a future record. If record was done, it couldn't be deleted!
// @Tags Records
// @Param   cancelReq body recordto.RecordCancelation true "cancel description"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/info/{recordID} [put]
// Удаление только ожидаемой записи, а не уже совершённой.
func (rec *RecordCtrl) RecordCancel(w http.ResponseWriter, r *http.Request) {
	req := &recordto.RecordCancelation{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := rec.usecase.RecordCancel(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
