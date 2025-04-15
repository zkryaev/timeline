package records

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/query"
	"timeline/internal/controller/scope"
	"timeline/internal/entity/dto/recordto"

	"go.uber.org/zap"
)

type Record interface {
	Record(ctx context.Context, logger *zap.Logger, param recordto.RecordParam) (*recordto.RecordList, error)
	RecordList(ctx context.Context, logger *zap.Logger, params *recordto.RecordListParams) (*recordto.RecordList, error)
	RecordAdd(ctx context.Context, logger *zap.Logger, rec *recordto.Record) error
	RecordCancel(ctx context.Context, logger *zap.Logger, rec *recordto.RecordCancelation) error
	Feedback
}

type RecordCtrl struct {
	usecase    Record
	Logger     *zap.Logger
	middleware middleware.Middleware
	settings   *scope.Settings
}

func New(usecase Record, middleware middleware.Middleware, logger *zap.Logger, settings *scope.Settings) *RecordCtrl {
	return &RecordCtrl{
		usecase:    usecase,
		Logger:     logger,
		middleware: middleware,
		settings:   settings,
	}
}

// @Summary Record information
// @Description `Если авторизация отключена: `user_id` или `org_id` нужно прокидывать в параметрах, иначе придет пустота - что логично`
// @Description (Если оба прокинуты, то выберется user_id - это влияет только на время полученных записей)
// @Description Типы Required параметров
// @Description `as_list=false` - (ОБЯЗАТЕЛЕН: record_id) возвращает данные одной записи.
// @Description  `as_list=true` - (ОБЯЗАТЕЛЕН: limit, page) возвращает список записей с пагинацией
// @Tags records
// @Param record_id query int false " "
// @Param limit query int false " "
// @Param page query int false " "
// @Param user_id query int false " "
// @Param org_id query int false " "
// @Param as_list query bool false " "
// @Param fresh query bool false "true - сегодняшние и будущие записи. false/not_given - до текущего дня"
// @Success 200 {object} recordto.RecordScrap "as_list=false"
// @Success 200 {object} recordto.RecordList "as_list=true"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /records [get]
func (rec *RecordCtrl) Record(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(rec.settings, rec.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(rec.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		asList = query.NewParamBool(scope.AS_LIST, false)
		orgID  = query.NewParamInt(scope.ORG_ID, false)
		userID = query.NewParamInt(scope.USER_ID, false)
	)
	params := query.NewParams(rec.settings, asList)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !rec.settings.EnableAuthorization {
		if orgID.Val != 0 {
			tdata.ID = orgID.Val
			tdata.IsOrg = true
		}
		if userID.Val != 0 {
			tdata.ID = userID.Val
			tdata.IsOrg = false
		}
	}
	var data *recordto.RecordList
	switch asList.Val {
	case scope.LIST:
		var (
			limit = query.NewParamInt(scope.LIMIT, true)
			page  = query.NewParamInt(scope.PAGE, true)
			fresh = query.NewParamBool(scope.FRESH, false)
		)
		params = query.NewParams(rec.settings, orgID, userID, limit, page, fresh)
		if err := params.Parse(r.URL.Query()); err != nil {
			logger.Error("param.Parse", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		req := &recordto.RecordListParams{TData: tdata, OrgID: orgID.Val, UserID: userID.Val, Limit: limit.Val, Page: page.Val, Fresh: fresh.Val}
		data, err = rec.usecase.RecordList(r.Context(), logger, req)
	case scope.SINGLE:
		recordID := query.NewParamInt(scope.RECORD_ID, true)
		params = query.NewParams(rec.settings, recordID)
		if err := params.Parse(r.URL.Query()); err != nil {
			logger.Error("param.Parse", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		req := recordto.RecordParam{RecordID: recordID.Val, TData: tdata}
		data, err = rec.usecase.Record(r.Context(), logger, req)
	}
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Records", zap.Bool(scope.AS_LIST, asList.Val), zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Records", zap.Bool(scope.AS_LIST, asList.Val), zap.Error(err))
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

// @Summary Add record
// @Description
// @Tags records
// @Accept  json
// @Param record body recordto.Record true " "
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /records [post]
func (rec *RecordCtrl) RecordAdd(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(rec.settings, rec.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(rec.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &recordto.Record{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if rec.settings.EnableAuthorization {
		req.UserID = tdata.ID
	}
	if err := rec.usecase.RecordAdd(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("RecordList", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("RecordList", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a future record
// @Description `Если авторизация отключена: `user_id` или `org_id` в теле запроса надо прокидывать, иначе не удалится`
// @Description (Если прокинуты обе, будет использован user_id)
// @Description Удаление `ожидаемой` записи. Если запись выполнена, ее не получится удалить
// @Tags records
// @Param   req body recordto.RecordCancelation true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /records [put]
func (rec *RecordCtrl) RecordCancel(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(rec.settings, rec.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(rec.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &recordto.RecordCancelation{TData: tdata}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !rec.settings.EnableAuthorization {
		if req.OrgID != 0 {
			req.TData.ID = req.OrgID
			req.TData.IsOrg = true
		} else {
			req.TData.ID = req.UserID
			req.TData.IsOrg = false
		}
	}
	if err := rec.usecase.RecordCancel(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("RecordCancel", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("RecordCancel", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
