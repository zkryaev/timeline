package orgs

import (
	"context"
	"net/http"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
)

type Services interface {
	Service(ctx context.Context, ServiceID, OrgID int) (*orgdto.ServiceResp, error)
	ServiceWorkerList(ctx context.Context, ServiceID, OrgID int) ([]*orgdto.WorkerResp, error)
	ServiceAdd(ctx context.Context, Service *orgdto.AddServiceReq) (*orgdto.ServiceResp, error)
	ServiceUpdate(ctx context.Context, Service *orgdto.UpdateServiceReq) (*orgdto.UpdateServiceReq, error)
	ServiceList(ctx context.Context, OrgID int) ([]*orgdto.ServiceResp, error)
	ServiceDelete(ctx context.Context, ServiceID, OrgID int) error
}

// @Summary Get service
// @Description Get specified service for specified organization
// @Tags Organization/Services
// @Produce json
// @Param   orgID path int true "org_id"
// @Param   serviceID path int true "service_id"
// @Success 200 {object} orgdto.ServiceResp
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services/{serviceID} [get]
func (o *OrgCtrl) Service(w http.ResponseWriter, r *http.Request) {
	path, err := validation.FetchSpecifiedID(mux.Vars(r), "orgID", "serviceID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := o.usecase.Service(ctx, path["serviceID"], path["orgID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary Get service workers
// @Description Get all workers that perform specified service in specified organization
// @Tags Organization/Services
// @Produce json
// @Param   orgID path int true "org_id"
// @Param   serviceID path int true "service_id"
// @Success 200 {array} orgdto.WorkerResp
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services/{serviceID}/workers [get]
func (o *OrgCtrl) ServiceWorkerList(w http.ResponseWriter, r *http.Request) {
	path, err := validation.FetchSpecifiedID(mux.Vars(r), "orgID", "serviceID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := o.usecase.ServiceWorkerList(ctx, path["serviceID"], path["orgID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary Add service
// @Description Add service for specified organization
// @Tags Organization/Services
// @Accept json
// @Produce json
// @Param   request body orgdto.AddServiceReq true "service info"
// @Success 200 {object} orgdto.ServiceResp
// @Failure 400
// @Failure 500
// @Router /orgs/services [post]
func (o *OrgCtrl) ServiceAdd(w http.ResponseWriter, r *http.Request) {
	var req orgdto.AddServiceReq
	if o.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	var err error
	if err = o.validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := o.usecase.ServiceAdd(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary Update service
// @Description Update specified service for specified organization
// @Tags Organization/Services
// @Accept json
// @Produce json
// @Param   request body orgdto.UpdateServiceReq true "service info"
// @Success 200 {object} orgdto.UpdateServiceReq
// @Failure 400
// @Failure 500
// @Router /orgs/services [put]
func (o *OrgCtrl) ServiceUpdate(w http.ResponseWriter, r *http.Request) {
	var req orgdto.UpdateServiceReq
	if o.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	var err error
	if err = o.validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := o.usecase.ServiceUpdate(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary List services
// @Description Get list of services for specified organization
// @Tags Organization/Services
// @Produce json
// @Param   orgID path int true "org_id"
// @Success 200 {array} orgdto.ServiceResp
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services [get]
func (o *OrgCtrl) ServiceList(w http.ResponseWriter, r *http.Request) {
	path, err := validation.FetchSpecifiedID(mux.Vars(r), "orgID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := o.usecase.ServiceList(ctx, path["orgID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary Delete service
// @Description Delete specified service for specified organization
// @Tags Organization/Services
// @Param   orgID path int true "org_id"
// @Param   serviceID path int true "service_id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services/{serviceID} [delete]
func (o *OrgCtrl) ServiceDelete(w http.ResponseWriter, r *http.Request) {
	path, err := validation.FetchSpecifiedID(mux.Vars(r), "orgID", "serviceID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if err := o.usecase.ServiceDelete(ctx, path["serviceID"], path["orgID"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
