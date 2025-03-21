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
	Organization(ctx context.Context, id int) (*orgdto.Organization, error)
	OrgUpdate(ctx context.Context, org *orgdto.OrgUpdateReq) error
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
	params, err := validation.FetchPathID(mux.Vars(r), "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Organization(r.Context(), params["id"])
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}
	if common.WriteJSON(w, data) != nil {
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
	req := &orgdto.OrgUpdateReq{}
	if common.DecodeAndValidate(r, req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.OrgUpdate(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
