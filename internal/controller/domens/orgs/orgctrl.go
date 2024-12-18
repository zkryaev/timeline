package orgs

import (
	"context"
	"net/http"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
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
	usecase   Org
	Logger    *zap.Logger
	json      jsoniter.API
	validator *validator.Validate
}

func NewOrgCtrl(usecase Org, logger *zap.Logger, jsoniter jsoniter.API, validator *validator.Validate) *OrgCtrl {
	return &OrgCtrl{
		usecase:   usecase,
		Logger:    logger,
		json:      jsoniter,
		validator: validator,
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
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
	if o.json.NewDecoder(r.Body).Decode(req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if err := o.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := o.usecase.OrgUpdate(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
