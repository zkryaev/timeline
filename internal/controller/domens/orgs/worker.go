package orgs

import (
	"context"
	"net/http"
	"strconv"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
)

type Workers interface {
	Worker(ctx context.Context, workerID, OrgID int) (*orgdto.WorkerResp, error)
	WorkerAdd(ctx context.Context, worker *orgdto.AddWorkerReq) (*orgdto.WorkerResp, error)
	WorkerUpdate(ctx context.Context, worker *orgdto.UpdateWorkerReq) error
	WorkerAssignService(ctx context.Context, assignInfo *orgdto.AssignWorkerReq) error
	WorkerUnAssignService(ctx context.Context, assignInfo *orgdto.AssignWorkerReq) error
	WorkerList(ctx context.Context, OrgID, Limit, Page int) (*orgdto.WorkerList, error)
	WorkerDelete(ctx context.Context, WorkerID, OrgID int) error
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
	path, err := validation.FetchPathID(mux.Vars(r), "orgID", "workerID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Worker(r.Context(), path["workerID"], path["orgID"])
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if common.WriteJSON(w, data) != nil {
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
	req := &orgdto.AddWorkerReq{}
	if common.DecodeAndValidate(r, req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.WorkerAdd(r.Context(), req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if common.WriteJSON(w, data) != nil {
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
	req := &orgdto.UpdateWorkerReq{}
	if common.DecodeAndValidate(r, req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.WorkerUpdate(r.Context(), req); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
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
	req := &orgdto.AssignWorkerReq{}
	if common.DecodeAndValidate(r, req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.WorkerAssignService(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
	params, err := validation.FetchPathID(mux.Vars(r), "orgID", "workerID", "serviceID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &orgdto.AssignWorkerReq{
		ServiceID: params["serviceID"],
		OrgID:     params["orgID"],
		WorkerID:  params["workerID"],
	}
	if err = o.usecase.WorkerUnAssignService(r.Context(), req); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
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
	path, err := validation.FetchPathID(mux.Vars(r), "orgID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := map[string]bool{
		"limit": true,
		"page":  true,
	}
	if !validation.IsQueryValid(r, query) {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	data, err := o.usecase.WorkerList(r.Context(), path["orgID"], limit, page)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if common.WriteJSON(w, data) != nil {
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
	path, err := validation.FetchPathID(mux.Vars(r), "orgID", "workerID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if o.usecase.WorkerDelete(r.Context(), path["workerID"], path["orgID"]) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
