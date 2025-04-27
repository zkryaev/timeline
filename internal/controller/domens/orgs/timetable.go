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

type Timetable interface {
	Timetable(ctx context.Context, logger *zap.Logger, req orgdto.TimetableReq) (*orgdto.Timetable, error)
	TimetableAdd(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error
	TimetableUpdate(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error
	TimetableDelete(ctx context.Context, logger *zap.Logger, orgID, weekday int) error
}

// @Summary Get timetable
// @Description Get organization timetable
// @Description `Если авторизация отключена, то время будет в часовом поясе организации из параметров`
// @Description
// @Description If user made call THEN org_id - mustbe
// @Description If org made call THEN org_id = token ID
// @Tags orgs/timetables
// @Accept  json
// @Param   org_id query int true " "
// @Success 200 {object} orgdto.Timetable
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/timetables [get]
func (o *OrgCtrl) Timetable(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Error("TokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	orgID := &query.IntParam{}
	if !tdata.IsOrg && !o.settings.EnableAuthorization {
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
	params := query.NewParams(o.settings, orgID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !o.settings.EnableAuthorization {
		tdata.ID = orgID.Val
		tdata.IsOrg = true
	}
	req := orgdto.TimetableReq{OrgID: orgID.Val, TData: tdata}
	data, err := o.usecase.Timetable(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Error("Timetable", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Timetable", zap.Error(err))
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

// @Summary Add timetable
// @Description Добавить расписание для организации
// @Description `Если авторизация отключена: прокидывать `org_id` в теле запроса`
// @Tags orgs/timetables
// @Accept  json
// @Param   request body orgdto.Timetable true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/timetables [post]
func (o *OrgCtrl) TimetableAdd(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.Timetable{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
	}
	if err := o.usecase.TimetableAdd(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrTimeIncorrect):
			logger.Info("TimetableAdd", zap.Error(err))
			http.Error(w, common.ErrTimeIncorrect.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("TimetableAdd", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("TimetableAdd", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update timetable
// @Description Обновление расписания организации
// @Description `Если авторизация отключена: прокидывать `org_id` в теле запроса`
// @Tags orgs/timetables
// @Accept  json
// @Param   request body orgdto.Timetable true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/timetable [put]
func (o *OrgCtrl) TimetableUpdate(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.Timetable{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
	}
	if err := o.usecase.TimetableUpdate(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrTimeIncorrect):
			logger.Info("TimetableUpdate", zap.Error(err))
			http.Error(w, common.ErrTimeIncorrect.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("TimetableUpdate", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("TimetableUpdate", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete timetable
// @Description Удаление расписания организации, если `weekday` выставлен, то удалится заданный день
// @Description `Если авторизация отключена: прокидывать `org_id` в параметрах`
// @Tags orgs/timetables
// @Accept  json
// @Param org_id query int false " "
// @Param weekday query int false " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/timetables [delete]
func (o *OrgCtrl) TimetableDelete(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		orgID   = query.NewParamInt(scope.ORG_ID, false)
		weekday = query.NewParamInt(scope.WEEKDAY, false)
	)
	params := query.NewParams(o.settings, weekday, orgID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !o.settings.EnableAuthorization {
		tdata.ID = orgID.Val
	}
	if o.usecase.TimetableDelete(r.Context(), logger, tdata.ID, weekday.Val) != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("TimetableDelete", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("TimetableDelete", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
