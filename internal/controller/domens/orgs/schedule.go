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

type Schedule interface {
	WorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) (*orgdto.ScheduleList, error)
	AddWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error
	UpdateWorkerSchedule(ctx context.Context, logger *zap.Logger, schedule *orgdto.WorkerSchedule) error
	DeleteWorkerSchedule(ctx context.Context, logger *zap.Logger, params *orgdto.ScheduleParams) error
}

// @Summary Get worker schedule
// @Description Get specified org worker's schedule with weekday filter (also may provide worker_id to narrow result). If no weekday then all week will be returned
// @Description
// @Description If user made call THEN org_id - mustbe
// @Description If org made call THEN org_id = token ID
// @Tags orgs/workers/schedule
// @Produce json
// @Param orgID query int true " "
// @Param workerID query int true " "
// @Param weekday query int false " "
// @Param limit query int true " "
// @Param page query int true " "
// @Success 200 {object} orgdto.ScheduleList
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/workers/schedules [get]
func (o *OrgCtrl) WorkersSchedule(w http.ResponseWriter, r *http.Request) {
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
	var (
		workerID = query.NewParamInt(scope.WORKER_ID, false)
		weekday  = query.NewParamInt(scope.WEEKDAY, false)
		limit    = query.NewParamInt(scope.LIMIT, true)
		page     = query.NewParamInt(scope.PAGE, true)
	)
	params := query.NewParams(o.settings, orgID, workerID, weekday, limit, page)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &orgdto.ScheduleParams{
		OrgID:    orgID.Val,
		WorkerID: workerID.Val,
		Weekday:  weekday.Val,
		Limit:    limit.Val,
		Page:     page.Val,
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
// @Description Удаление расписания работника организации. Если указан weekday, то удален будет только заданный день
// @Tags orgs/workers/schedule
// @Param   worker_id query int true " "
// @Param   weekday query int false " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers/schedules [delete]
func (o *OrgCtrl) DeleteWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		workerID = query.NewParamInt(scope.WORKER_ID, true)
		weekday  = query.NewParamInt(scope.WEEKDAY, false)
	)
	params := query.NewParams(o.settings, workerID, weekday)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &orgdto.ScheduleParams{
		OrgID:    scope.DEAD_ORG_ID,
		WorkerID: workerID.Val,
		Weekday:  weekday.Val,
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
// @Tags orgs/workers/schedule
// @Accept json
// @Param   schedule body orgdto.WorkerSchedule true " "
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers/schedules [put]
func (o *OrgCtrl) UpdateWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.WorkerSchedule{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
// @Tags orgs/workers/schedule
// @Accept json
// @Produce json
// @Param   schedule body orgdto.WorkerSchedule true "Schedule data"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/workers/schedules [post]
func (o *OrgCtrl) AddWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.WorkerSchedule{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = tdata.ID
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
