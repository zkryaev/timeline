package domens

import (
	"context"
	"net/http"
	"timeline/internal/entity/dto"

	"github.com/go-playground/validator"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type User interface {
	SearchOrgs(ctx context.Context, sreq *dto.SearchReq) (*dto.SearchResp, error)
	OrgsInArea(ctx context.Context, area *dto.OrgAreaReq) (*dto.OrgAreaResp, error)
}

type UserCtrl struct {
	usecase   User
	Logger    *zap.Logger
	json      jsoniter.API
	validator validator.Validate
}

func NewUserCtrl(usecase User, logger *zap.Logger, jsoniter jsoniter.API, validator validator.Validate) *UserCtrl {
	return &UserCtrl{
		usecase:   usecase,
		Logger:    logger,
		json:      jsoniter,
		validator: validator,
	}
}

func (u *UserCtrl) SearchOrganization(w http.ResponseWriter, r *http.Request) {
	var req dto.SearchReq
	if u.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	// валидация полей
	if err := u.validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := u.usecase.SearchOrgs(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// отдаем токен
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if u.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

func (u *UserCtrl) OrganizationInArea(w http.ResponseWriter, r *http.Request) {
	var req dto.OrgAreaReq
	if u.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	// валидация полей
	if err := u.validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := u.usecase.OrgsInArea(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// отдаем токен
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if u.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}
