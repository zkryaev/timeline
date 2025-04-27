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

type Org interface {
	Organization(ctx context.Context, logger *zap.Logger, id int) (*orgdto.Organization, error)
	OrgUpdate(ctx context.Context, logger *zap.Logger, org *orgdto.OrgUpdateReq) error
	Timetable
	Workers
	Services
	Slots
	Schedule
}

type OrgCtrl struct {
	usecase    Org
	Logger     *zap.Logger
	middleware middleware.Middleware
	settings   *scope.Settings
}

func New(usecase Org, middleware middleware.Middleware, logger *zap.Logger, settings *scope.Settings) *OrgCtrl {
	return &OrgCtrl{
		usecase:    usecase,
		Logger:     logger,
		middleware: middleware,
		settings:   settings,
	}
}

// @Summary Full organization info
// @Description
// @Tags orgs
// @Param   org_id query int true "database id"
// @Description
// @Description If user made call THEN org_id - mustbe
// @Description If org made call THEN org_id = token ID
// @Success 200 {object} orgdto.Organization
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs [get]
func (o *OrgCtrl) GetOrganization(w http.ResponseWriter, r *http.Request) {
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
	data, err := o.usecase.Organization(r.Context(), logger, orgID.Val)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Organization", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Organization", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	if common.WriteJSON(w, data) != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Change organization info
// @Description
// @Tags orgs
// @Accept  json
// @Param   request body orgdto.OrgUpdateReq true " "
// @Success 200
// @Failure 400
// @Failure 304
// @Failure 500
// @Router /orgs [put]
func (o *OrgCtrl) PutOrganization(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(o.settings, o.Logger, r.Context())
	token, err := middleware.GetTokenDataFromCtx(o.settings, r.Context())
	if err != nil {
		logger.Error("TokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &orgdto.OrgUpdateReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.settings.EnableAuthorization {
		req.OrgID = token.ID
	}
	if err := o.usecase.OrgUpdate(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("OrgUpdate", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("OrgUpdate", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
