package orgs

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Workers interface {
	Worker(ctx context.Context, logger *zap.Logger, workerID, OrgID int) (*orgdto.WorkerResp, error)
	WorkerAdd(ctx context.Context, logger *zap.Logger, worker *orgdto.AddWorkerReq) (*orgdto.WorkerResp, error)
	WorkerUpdate(ctx context.Context, logger *zap.Logger, worker *orgdto.UpdateWorkerReq) error
	WorkerAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error
	WorkerUnAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error
	WorkerList(ctx context.Context, logger *zap.Logger, OrgID, Limit, Page int) (*orgdto.WorkerList, error)
	WorkerDelete(ctx context.Context, logger *zap.Logger, WorkerID, OrgID int) error
}

// @Summary Get worker
// @Description Get specified worker for specified organization
// @Tags organization/workers
// @Produce json
// @Param   orgID path int true "org_id"
// @Param   workerID path int true "worker_id"
// @Success 200 {object} orgdto.WorkerResp
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/workers/{workerID} [get]
func (o *OrgCtrl) Worker(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	path, err := validation.FetchPathID(mux.Vars(r), "orgID", "workerID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Worker(r.Context(), logger, path["workerID"], path["orgID"])
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Worker", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Worker", zap.Error(err))
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

// @Summary Add worker
// @Description Add a new worker to the specified organization
// @Tags organization/workers
// @Accept json
// @Produce json
// @Param   request body orgdto.AddWorkerReq true "worker info"
// @Success 200 {object} orgdto.WorkerResp
// @Failure 400
// @Failure 500
// @Router /orgs/workers [post]
func (o *OrgCtrl) WorkerAdd(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.AddWorkerReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.WorkerAdd(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("WorkerAdd", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("WorkerAdd", zap.Error(err))
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

// @Summary Update worker
// @Description Update specified worker for specified organization
// @Tags organization/workers
// @Accept json
// @Param   request body orgdto.UpdateWorkerReq true "worker info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/workers [put]
func (o *OrgCtrl) WorkerUpdate(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.UpdateWorkerReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.WorkerUpdate(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("WorkerUpdate", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("WorkerUpdate", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Assign worker to service
// @Description Assign a specified worker to a specified service in the specified organization
// @Tags organization/workers
// @Accept json
// @Produce json
// @Param   request body orgdto.AssignWorkerReq true "assignment info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/workers/service [post]
func (o *OrgCtrl) WorkerAssignService(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.AssignWorkerReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.WorkerAssignService(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("WorkerAssignService", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("WorkerAssignService", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Unassign worker from
// @Description Unassign worker from specified organization service
// @Tags organization/workers
// @Accept json
// @Produce json
// @Param   orgID path int true "org_id"
// @Param   workerID path int true "worker_id"
// @Param   serviceID path int true "service_id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/workers/{workerID}/service/{serviceID} [delete]
func (o *OrgCtrl) WorkerUnAssignService(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "orgID", "workerID", "serviceID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &orgdto.AssignWorkerReq{
		ServiceID: params["serviceID"],
		OrgID:     params["orgID"],
		WorkerID:  params["workerID"],
	}
	if err = o.usecase.WorkerUnAssignService(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("WorkerUnAssignService", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("WorkerUnAssignService", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary List workers
// @Description Get list of workers for specified organization
// @Tags organization/workers
// @Produce json
// @Param   orgID path int true "org_id"
// @Param limit query int true "Limit the number of results"
// @Param page query int true "Page number for pagination"
// @Success 200 {array} orgdto.WorkerList
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/workers [get]
func (o *OrgCtrl) WorkerList(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	path, err := validation.FetchPathID(mux.Vars(r), "orgID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	query := map[string]bool{
		"limit": true,
		"page":  true,
	}
	if err := validation.IsQueryValid(r, query); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	data, err := o.usecase.WorkerList(r.Context(), logger, path["orgID"], limit, page)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("WorkerList", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("WorkerList", zap.Error(err))
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

// @Summary Delete worker
// @Description Delete specified worker from specified organization
// @Tags organization/workers
// @Param   orgID path int true "org_id"
// @Param   workerID path int true "worker_id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/workers/{workerID} [delete]
func (o *OrgCtrl) WorkerDelete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	path, err := validation.FetchPathID(mux.Vars(r), "orgID", "workerID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.WorkerDelete(r.Context(), logger, path["workerID"], path["orgID"]); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("WorkerDelete", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("WorkerDelete", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
