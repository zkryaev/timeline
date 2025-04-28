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

type Slots interface {
	Slots(ctx context.Context, logger *zap.Logger, req *orgdto.SlotReq) ([]*orgdto.SlotResp, error)
}

// @Summary Get slots
// @Description Получение всех слотов работника
// @Description `Если авторизация отключена, то время будет в часовом поясе организации из параметров`
// @Description
// @Description If user made call THEN org_id - mustbe
// @Description If org made call THEN org_id = token ID
// @Tags orgs/workers/slots
// @Produce json
// @Param   worker_id query int true " "
// @Param   org_id query int true " "
// @Success 200 {array} orgdto.SlotResp
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/workers/slots [get]
func (o *OrgCtrl) Slots(w http.ResponseWriter, r *http.Request) {
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
	workerID := query.NewParamInt(scope.WORKER_ID, true)
	params := query.NewParams(o.settings, orgID, workerID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !o.settings.EnableAuthorization {
		tdata.ID = orgID.Val
		tdata.IsOrg = true
	}
	req := &orgdto.SlotReq{OrgID: orgID.Val, WorkerID: workerID.Val, TData: tdata}
	data, err := o.usecase.Slots(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Slots", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Slots", zap.Error(err))
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
