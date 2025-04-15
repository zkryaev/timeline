package orgs

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/query"
	"timeline/internal/controller/scope"
	"timeline/internal/entity/dto/orgdto"

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

// @Summary Get organization's worker
// @Description
// Если `as_list=false` - (ОБЯЗАТЕЛЕН worker_id) возвращает данные одного работника.
// Если `as_list=true` -  (НЕТ) возвращает список работников с пагинацией
// @Tags orgs/workers
// @Produce json
// @Param org_id query int true " "
// @Param worker_id query int true " "
// @Param as_list query bool false " "
// @Param limit query int false " "
// @Param page query int false " "
// @Success 200 {object} orgdto.WorkerResp "as_list=false"
// @Success 200 {object} orgdto.WorkerList "as_list=true"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/workers [get]
func (o *OrgCtrl) Workers(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	var (
		orgID  = query.NewParamInt(scope.ORG_ID, true)
		asList = query.NewParamBool(scope.AS_LIST, false)
	)
	params := query.NewParams(o.settings, orgID, asList)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var data any
	var err error
	switch asList.Val {
	case scope.LIST:
		limit := query.NewParamInt(scope.LIMIT, true)
		page := query.NewParamInt(scope.PAGE, true)
		params = query.NewParams(o.settings, limit, page)
		if err := params.Parse(r.URL.Query()); err != nil {
			logger.Error("param.Parse", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		data, err = o.usecase.WorkerList(r.Context(), logger, orgID.Val, limit.Val, page.Val)
	case scope.SINGLE:
		workerID := query.NewParamInt(scope.WORKER_ID, true)
		params = query.NewParams(o.settings, workerID)
		if err := params.Parse(r.URL.Query()); err != nil {
			logger.Error("param.Parse", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		data, err = o.usecase.Worker(r.Context(), logger, workerID.Val, orgID.Val)
	}
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Workers", zap.Bool(scope.AS_LIST, asList.Val), zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Workers", zap.Bool(scope.AS_LIST, asList.Val), zap.Error(err))
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
// @Description
// @Tags orgs/workers
// @Accept json
// @Produce json
// @Param   request body orgdto.AddWorkerReq true " "
// @Success 200 {object} orgdto.WorkerResp
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers [post]
func (o *OrgCtrl) WorkerAdd(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.AddWorkerReq{
		OrgID: tdata.ID,
	}
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

// @Summary Update worker info
// @Description
// @Tags orgs/workers
// @Accept json
// @Param   request body orgdto.UpdateWorkerReq true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers [put]
func (o *OrgCtrl) WorkerUpdate(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.UpdateWorkerReq{OrgID: tdata.ID}
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

// @Summary Delete worker
// @Description Delete specified worker from specified organization
// @Tags orgs/workers
// @Param   worker_id query int true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers [delete]
func (o *OrgCtrl) WorkerDelete(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		workerID = query.NewParamInt(scope.WORKER_ID, true)
	)
	params := query.NewParams(o.settings, workerID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.WorkerDelete(r.Context(), logger, workerID.Val, tdata.ID); err != nil {
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

// @Summary Assign worker to service
// @Description Assign a specified worker to a specified service in the specified organization
// @Tags orgs/workers
// @Accept json
// @Produce json
// @Param   request body orgdto.AssignWorkerReq true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers/services [post]
func (o *OrgCtrl) WorkerAssignService(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.AssignWorkerReq{OrgID: tdata.ID}
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
// @Tags orgs/workers
// @Accept json
// @Produce json
// @Param   worker_id query int true " "
// @Param   service_id query int true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers/services [delete]
func (o *OrgCtrl) WorkerUnassignService(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		serviceID = query.NewParamInt(scope.SERVICE_ID, true)
		workerID  = query.NewParamInt(scope.WORKER_ID, true)
	)
	params := query.NewParams(o.settings, serviceID, workerID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &orgdto.AssignWorkerReq{
		OrgID:     tdata.ID,
		ServiceID: serviceID.Val,
		WorkerID:  workerID.Val,
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
