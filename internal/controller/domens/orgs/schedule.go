package orgs

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/sugar/custom"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Schedule interface {
	WorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) (*orgdto.ScheduleList, error)
	AddWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error
	UpdateWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error
	DeleteWorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) error
}

// @Summary Get worker schedule
// @Description Get specified worker schedule for specified org with weekday filter. If no weekday then all week will be returned
// @Tags organization/schedule
// @Produce json
// @Param   orgID path int true "org_id"
// @Param   workerID query int true "Returned schedule for specified worker, otherwise for all org's workers"
// @Param weekday query int false "weekday"
// @Param limit query int true "Limit the number of results"
// @Param page query int true "Page number for pagination"
// @Success 200 {object} orgdto.ScheduleList
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/{orgID}/schedules [get]
func (o *OrgCtrl) WorkerSchedule(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "orgID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	query := map[string]string{
		"worker_id": "int",
		"weekday":   "int",
		"limit":     "int",
		"page":      "int",
	}
	queryParams, err := custom.QueryParamsConv(query, r.URL.Query())
	if err != nil {
		logger.Error("QueryParamsConv", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &orgdto.ScheduleParams{
		WorkerID: queryParams["worker_id"].(int),
		OrgID:    params["orgID"],
		Weekday:  queryParams["weekday"].(int),
		Limit:    queryParams["limit"].(int),
		Page:     queryParams["page"].(int),
	}
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.WorkerSchedule(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("WorkerSchedule", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("WorkerSchedule", zap.Error(err))
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

// @Summary Delete worker schedule
// @Description Delete specified worker schedule for a specific organization and worker with an optional weekday filter
// @Tags organization/schedule
// @Param   workerID path int true "worker_id"
// @Param   orgID path int true "org_id"
// @Param   weekday query int false "weekday"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/schedules/{workerID} [delete]
func (o *OrgCtrl) DeleteWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "orgID", "workerID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var weekday int
	if r.URL.Query().Get("weekday") != "" {
		weekday, err = strconv.Atoi(r.URL.Query().Get("weekday"))
		if err != nil {
			logger.Error("Atoi", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	req := &orgdto.ScheduleParams{
		WorkerID: params["workerID"],
		OrgID:    params["orgID"],
		Weekday:  weekday,
	}
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.DeleteWorkerSchedule(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("DeleteWorkerSchedule", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("DeleteWorkerSchedule", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update worker schedule
// @Description Update the schedule for a specific worker in an organization
// @Tags organization/schedule
// @Accept json
// @Param   schedule body orgdto.WorkerSchedule true "Schedule data"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/schedules [put]
func (o *OrgCtrl) UpdateWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.WorkerSchedule{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.UpdateWorkerSchedule(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("UpdateWorkerSchedule", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("UpdateWorkerSchedule", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Add worker schedule
// @Description Add a new schedule for a specific worker in an organization
// @Tags organization/schedule
// @Accept json
// @Produce json
// @Param   schedule body orgdto.WorkerSchedule true "Schedule data"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/schedules [post]
func (o *OrgCtrl) AddWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.WorkerSchedule{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.AddWorkerSchedule(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrTimeIncorrect):
			logger.Info("AddWorkerSchedule", zap.Error(err))
			http.Error(w, common.ErrTimeIncorrect.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("AddWorkerSchedule", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("AddWorkerSchedule", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
