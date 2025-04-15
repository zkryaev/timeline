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

type Services interface {
	Service(ctx context.Context, logger *zap.Logger, ServiceID, OrgID int) (*orgdto.ServiceList, error)
	ServiceList(ctx context.Context, logger *zap.Logger, OrgID, Limit, Page int) (*orgdto.ServiceList, error)
	ServiceAdd(ctx context.Context, logger *zap.Logger, Service *orgdto.AddServiceReq) error
	ServiceUpdate(ctx context.Context, logger *zap.Logger, Service *orgdto.UpdateServiceReq) error
	ServiceDelete(ctx context.Context, logger *zap.Logger, ServiceID, OrgID int) error
}

// @Summary Get service
// @Description Типы Required параметров
// @Description `org_id` - всегда обязателен
// @Description Если `as_list=false` - (ОБЯЗАТЕЛЕН: service_id) возвращает данные одной услуги.
// @Description Если `as_list=true` -  (ОБЯЗАТЕЛЕН: limit, page) возвращает список услуг с пагинацией
// @Tags orgs/services
// @Produce json
// @Param org_id query int true " "
// @Param service_id query int false " "
// @Param limit query int false " "
// @Param page query int false " "
// @Param as_list query bool false " "
// @Success 200 {object} orgdto.ServiceList "as_list=true"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/services [get]
func (o *OrgCtrl) Service(w http.ResponseWriter, r *http.Request) {
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
	var data *orgdto.ServiceList
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
		data, err = o.usecase.ServiceList(r.Context(), logger, orgID.Val, limit.Val, page.Val)
	case scope.SINGLE:
		serviceID := query.NewParamInt(scope.SERVICE_ID, true)
		params = query.NewParams(o.settings, serviceID)
		if err := params.Parse(r.URL.Query()); err != nil {
			logger.Error("param.Parse", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		data, err = o.usecase.Service(r.Context(), logger, serviceID.Val, orgID.Val)
	}
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Services", zap.Bool(scope.AS_LIST, asList.Val), zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Services", zap.Bool(scope.AS_LIST, asList.Val), zap.Error(err))
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

// @Summary Add service
// @Description Добавление услуги
// @Description `Если авторизация отключена: `org_id`  прокинуть в тело запроса!`
// @Tags orgs/services
// @Accept json
// @Produce json
// @Param   request body orgdto.AddServiceReq true " "
// @Success 200 {object} orgdto.ServiceResp
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/services [post]
func (o *OrgCtrl) ServiceAdd(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.AddServiceReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
// @Tags orgs/services
// @Accept json
// @Param   request body orgdto.UpdateServiceReq true "service info"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/services [put]
func (o *OrgCtrl) ServiceUpdate(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.UpdateServiceReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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

// @Summary Delete service
// @Description Delete specified service for specified organization
// @Tags orgs/services
// @Param   service_id query int true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/services [delete]
func (o *OrgCtrl) ServiceDelete(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		serviceID = query.NewParamInt(scope.SERVICE_ID, true)
	)
	params := query.NewParams(o.settings, serviceID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	orgID := scope.DEAD_ORG_ID
	if o.settings.EnableAuthorization {
		orgID = tdata.ID
	}
	if err = o.usecase.ServiceDelete(r.Context(), logger, serviceID.Val, orgID); err != nil {
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
