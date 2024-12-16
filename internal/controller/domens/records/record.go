package records

import (
	"context"
	"net/http"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/recordto"
	"timeline/internal/libs/custom"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Record interface {
	Record(ctx context.Context, recordID int) (*recordto.RecordScrap, error)
	RecordList(ctx context.Context, params *recordto.RecordListParams) (*recordto.RecordList, error)
	RecordAdd(ctx context.Context, rec *recordto.Record) error
	RecordDelete(ctx context.Context, rec *recordto.Record) error
	Feedback
}

type RecordCtrl struct {
	usecase   Record
	Logger    *zap.Logger
	json      jsoniter.API
	validator *validator.Validate
}

func NewRecordCtrl(usecase Record, logger *zap.Logger, jsoniter jsoniter.API, validator *validator.Validate) *RecordCtrl {
	return &RecordCtrl{
		usecase:   usecase,
		Logger:    logger,
		json:      jsoniter,
		validator: validator,
	}
}

// @Summary Record information
// @Description Get bounded with record information
// @Tags Records
// @Param   recordID path int true "record_id"
// @Success 200 {object} recordto.RecordScrap
// @Failure 400
// @Failure 500
// @Router /records/{recordID} [get]
func (rec *RecordCtrl) Record(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "recordID")
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if rec.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
}

// @Summary Records
// @Description Get bounded with records informations
// @Tags Records
// @Param   userID query int true "user_id"
// @Param   orgID query int true "org_id"
// @Param   fresh query bool false "Decide which records must be returned. True - only current & future records. False/NotGiven - olds"
// @Success 200 {object} recordto.RecordList
// @Failure 400
// @Failure 500
// @Router /records/list [get]
func (rec *RecordCtrl) RecordList(w http.ResponseWriter, r *http.Request) {
	query := map[string]bool{
		"user_id": false,
		"org_id":  false,
		"fresh":  false,
	}
	if !validation.IsQueryValid(r, query) {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}
	params := map[string]string{
		"user_id": "int",
		"org_id":  "int",
		"fresh":  "bool",
	}
	queryParams, err := custom.QueryParamsConv(params, r.URL.Query())
	if err != nil {
		http.Error(w, "Invalid query parameters"+err.Error(), http.StatusBadRequest)
		return
	}
	req := &recordto.RecordListParams{
		UserID: queryParams["user_id"].(int),
		OrgID:  queryParams["org_id"].(int),
		Fresh:  queryParams["fresh"].(bool),
	}
	if err := rec.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := rec.usecase.RecordList(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if rec.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
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
	if rec.json.NewDecoder(r.Body).Decode(req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if err := rec.validator.Struct(req); err != nil {
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
// @Description Delete a future record. If record was done, it won't be deleted!
// @Tags Records
// @Param   recordID path int true "record_id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/{recordID} [delete]
// Удаление только ожидаемой записи, а не уже совершённой.
func (rec *RecordCtrl) RecordDelete(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "recordID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if params["recordID"] <= 0 {
		http.Error(w, "record_id must be > 0", http.StatusBadRequest)
		return
	}
	req := &recordto.Record{
		RecordID: params["recordID"],
	}
	if err := rec.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := rec.usecase.RecordDelete(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
