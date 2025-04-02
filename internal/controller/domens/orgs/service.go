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

type Services interface {
	Service(ctx context.Context, logger *zap.Logger, ServiceID, OrgID int) (*orgdto.ServiceResp, error)
	ServiceWorkerList(ctx context.Context, logger *zap.Logger, ServiceID, OrgID int) ([]*orgdto.WorkerResp, error)
	ServiceAdd(ctx context.Context, logger *zap.Logger, Service *orgdto.AddServiceReq) error
	ServiceUpdate(ctx context.Context, logger *zap.Logger, Service *orgdto.UpdateServiceReq) error
	ServiceList(ctx context.Context, logger *zap.Logger, OrgID, Limit, Page int) (*orgdto.ServiceList, error)
	ServiceDelete(ctx context.Context, logger *zap.Logger, ServiceID, OrgID int) error
}

// @Summary Get service
// @Description Get specified service for specified organization
// @Tags organization/services
// @Produce json
// @Param   orgID path int true "org_id"
// @Param   serviceID path int true "service_id"
// @Success 200 {object} orgdto.ServiceResp
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services/{serviceID} [get]
func (o *OrgCtrl) Service(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	path, err := validation.FetchPathID(mux.Vars(r), "orgID", "serviceID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Service(r.Context(), logger, path["serviceID"], path["orgID"])
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Service", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Service", zap.Error(err))
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

// @Summary Get service workers
// @Description Get all workers that perform specified service in specified organization
// @Tags organization/services
// @Produce json
// @Param   orgID path int true "org_id"
// @Param   serviceID path int true "service_id"
// @Success 200 {array} orgdto.WorkerResp
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services/{serviceID}/workers [get]
func (o *OrgCtrl) ServiceWorkerList(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	path, err := validation.FetchPathID(mux.Vars(r), "orgID", "serviceID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.ServiceWorkerList(r.Context(), logger, path["serviceID"], path["orgID"])
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("ServiceWorkerList", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("ServiceWorkerList", zap.Error(err))
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

// @Summary Add service
// @Description Add service for specified organization
// @Tags organization/services
// @Accept json
// @Produce json
// @Param   request body orgdto.AddServiceReq true "service info"
// @Success 200 {object} orgdto.ServiceResp
// @Failure 400
// @Failure 500
// @Router /orgs/services [post]
func (o *OrgCtrl) ServiceAdd(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.AddServiceReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.ServiceAdd(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("ServiceAdd", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("ServiceAdd", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update service
// @Description Update specified service for specified organization
// @Tags organization/services
// @Accept json
// @Param   request body orgdto.UpdateServiceReq true "service info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/services [put]
func (o *OrgCtrl) ServiceUpdate(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.UpdateServiceReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.ServiceUpdate(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("ServiceUpdate", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("ServiceUpdate", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary List services
// @Description Get list of services for specified organization
// @Tags organization/services
// @Produce json
// @Param   orgID path int true "org_id"
// @Param limit query int true "Limit the number of results"
// @Param page query int true "Page number for pagination"
// @Success 200 {array} orgdto.ServiceList
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services [get]
func (o *OrgCtrl) ServiceList(w http.ResponseWriter, r *http.Request) {
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
	data, err := o.usecase.ServiceList(r.Context(), logger, path["orgID"], limit, page)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("ServiceList", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("ServiceList", zap.Error(err))
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

// @Summary Delete service
// @Description Delete specified service for specified organization
// @Tags organization/services
// @Param   orgID path int true "org_id"
// @Param   serviceID path int true "service_id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/services/{serviceID} [delete]
func (o *OrgCtrl) ServiceDelete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	path, err := validation.FetchPathID(mux.Vars(r), "orgID", "serviceID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err = o.usecase.ServiceDelete(r.Context(), logger, path["serviceID"], path["orgID"]); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("ServiceDelete", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("ServiceDelete", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
