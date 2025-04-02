package orgs

import (
	"context"
	"net/http"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
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
	usecase Org
	Logger  *zap.Logger
}

func New(usecase Org, logger *zap.Logger) *OrgCtrl {
	return &OrgCtrl{
		usecase: usecase,
		Logger:  logger,
	}
}

// @Summary Organization information
// @Description Get organization information
// @Tags Organization
// @Param   id path int true "org_id"
// @Success 200 {object} orgdto.Organization
// @Failure 400
// @Failure 500
// @Router /orgs/info/{id} [get]
func (o *OrgCtrl) GetOrgByID(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "id")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Organization(r.Context(), logger, params["id"])
	if err != nil {
		logger.Error("Organization", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
	}
	if common.WriteJSON(w, data) != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Update org info
// @Description Update organization information
// @Tags Organization
// @Accept  json
// @Param   request body orgdto.OrgUpdateReq true "New org info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/update [put]
func (o *OrgCtrl) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.OrgUpdateReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.OrgUpdate(r.Context(), logger, req); err != nil {
		logger.Error("OrgUpdate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
