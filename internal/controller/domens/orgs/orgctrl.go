package orgs

import (
	"context"
	"net/http"
	"strconv"
	"timeline/internal/entity/dto/orgdto"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Org interface {
	Organization(ctx context.Context, id int) (*orgdto.Organization, error)
	OrgUpdate(ctx context.Context, org *orgdto.OrgUpdateReq) (*orgdto.OrgUpdateResp, error)
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
// @Router /orgs/info/{id} [put]
func (o *OrgCtrl) GetOrgByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idString, ok := params["id"]
	if !ok {
		http.Error(w, "No org id provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := o.usecase.Organization(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary Update
// @Description Update organization information
// @Tags Organization
// @Accept  json
// @Produce json
// @Param   request body orgdto.OrgUpdateReq true "New org info"
// @Success 200 {object} orgdto.OrgUpdateResp
// @Failure 400
// @Failure 500
// @Router /orgs/update [put]
func (o *OrgCtrl) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	var req orgdto.OrgUpdateReq
	if o.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	var err error
	if err = o.validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := o.usecase.OrgUpdate(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}
