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
	Worker(ctx context.Context, logger *zap.Logger, workerID, OrgID int) (*orgdto.WorkerList, error)
	WorkerAdd(ctx context.Context, logger *zap.Logger, worker *orgdto.AddWorkerReq) (*orgdto.WorkerResp, error)
	WorkerUpdate(ctx context.Context, logger *zap.Logger, worker *orgdto.UpdateWorkerReq) error
	WorkerList(ctx context.Context, logger *zap.Logger, OrgID, Limit, Page int) (*orgdto.WorkerList, error)
	WorkerDelete(ctx context.Context, logger *zap.Logger, WorkerID, OrgID int) error
	WorkerAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error
	WorkerUnAssignService(ctx context.Context, logger *zap.Logger, assignInfo *orgdto.AssignWorkerReq) error
	WorkersServices(ctx context.Context, logger *zap.Logger, ServiceID, OrgID int) ([]*orgdto.WorkerResp, error)
}

// @Summary Get organization's worker
// @Description Типы Required параметров
// @Description `org_id` - всегда обязателен
// @Description Если `as_list=false` - (ОБЯЗАТЕЛЕН:  worker_id) возвращает данные одного работника.
// @Description Если `as_list=true` -  (ОБЯЗАТЕЛЕН: limit, page) возвращает список работников с пагинацией
// @Description
// @Description If user made call THEN org_id - mustbe
// @Description If org made call THEN org_id = token ID
// @Tags orgs/workers
// @Produce json
// @Param org_id query int true " "
// @Param worker_id query int false " "
// @Param limit query int false " "
// @Param page query int false " "
// @Param as_list query bool false " "
// @Success 200 {object} orgdto.WorkerList "as_list=true"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/workers [get]
func (o *OrgCtrl) Workers(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Error("TokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	orgID := &query.IntParam{}
	if !tdata.IsOrg || !o.settings.EnableAuthorization {
		orgID = query.NewParamInt(scope.ORG_ID, true)
		params := query.NewParams(o.settings, orgID)
		if err := params.Parse(r.URL.Query()); err != nil {
			logger.Error("param.Parse", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	} else {
		orgID.Val = tdata.ID
	}
	asList := query.NewParamBool(scope.AS_LIST, false)
	params := query.NewParams(o.settings, orgID, asList)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var data *orgdto.WorkerList
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
// @Description Добавление работника к организации
// @Description `Если авторизация отключена: `org_id`  прокинуть в тело запроса!`
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
	req := &orgdto.AddWorkerReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
// @Description `Если авторизация отключена: `org_id` прокинуть в тело запроса!`
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
	req := &orgdto.UpdateWorkerReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
// @Description Удаление работника из организации
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
	if o.settings.EnableAuthorization {
		tdata.ID = scope.DEAD_ORG_ID
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

// @Summary Get service workers
// @Description Получение работников которые выполняют заданную услугу
// @Tags orgs/workers/services
// @Produce json
// @Param   org_id query int true " "
// @Param   service_id query int true " "
// @Success 200 {array} orgdto.WorkerResp
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/workers/services [get]
func (o *OrgCtrl) WorkersServices(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	var (
		orgID     = query.NewParamInt(scope.ORG_ID, true)
		serviceID = query.NewParamInt(scope.SERVICE_ID, true)
	)
	params := query.NewParams(o.settings, orgID, serviceID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.WorkersServices(r.Context(), logger, serviceID.Val, orgID.Val)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("WorkersServices", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("WorkersServices", zap.Error(err))
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

// @Summary Assign worker to service
// @Description Прикрепление заданного работника на выполнение заданной услуги организации
// @Description `Если авторизация отключена: `org_id` прокинуть в тело запроса!`
// @Tags orgs/workers/services
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
	req := &orgdto.AssignWorkerReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
// @Description Открепление работника от заданной услуги организации
// @Tags orgs/workers/services
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
		OrgID:     scope.DEAD_ORG_ID,
		ServiceID: serviceID.Val,
		WorkerID:  workerID.Val,
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
