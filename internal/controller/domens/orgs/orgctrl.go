package orgs

import (
	"context"
	"net/http"
	"timeline/internal/entity/dto/orgdto"

	"github.com/go-playground/validator"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Org interface {
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
