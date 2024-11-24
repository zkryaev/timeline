package users

import (
	"context"
	"net/http"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/entity/dto/userdto"

	"github.com/go-playground/validator"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type User interface {
	SearchOrgs(ctx context.Context, sreq *orgdto.SearchReq) (*orgdto.SearchResp, error)
	OrgsInArea(ctx context.Context, area *orgdto.OrgAreaReq) (*orgdto.OrgAreaResp, error)
	UserUpdate(ctx context.Context, user *userdto.UserUpdateReq) (*userdto.UserUpdateResp, error)
}

type UserCtrl struct {
	usecase   User
	Logger    *zap.Logger
	json      jsoniter.API
	validator *validator.Validate
}

func NewUserCtrl(usecase User, logger *zap.Logger, jsoniter jsoniter.API, validator *validator.Validate) *UserCtrl {
	return &UserCtrl{
		usecase:   usecase,
		Logger:    logger,
		json:      jsoniter,
		validator: validator,
	}
}

// @Summary Update
// @Description Update user information
// @Tags User
// @Accept  json
// @Produce  json
// @Param   request body userdto.UserUpdateReq true "New user info"
// @Success 200 {object} userdto.UserUpdateResp
// @Failure 400
// @Failure 500
// @Router /user/update [put]
func (u *UserCtrl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req userdto.UserUpdateReq
	if u.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	var err error
	if err = u.validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := u.usecase.UserUpdate(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if u.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary Organization Searching
// @Description Get organizations that are satisfiered to search params
// @Tags User
// @Accept  json
// @Produce  json
// @Param   request body orgdto.SearchReq true "Search parameters request"
// @Success 200 {object} orgdto.SearchResp
// @Failure 400
// @Failure 500
// @Router /user/find/orgs [get]
func (u *UserCtrl) SearchOrganization(w http.ResponseWriter, r *http.Request) {
	var req orgdto.SearchReq
	if u.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	var err error
	if err = u.validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := u.usecase.SearchOrgs(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if u.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
}

// @Summary Show Organizations on Map
// @Description Get organizations that located at the given area
// @Tags User
// @Accept  json
// @Produce  json
// @Param   request body orgdto.OrgAreaReq true "Area parameters"
// @Success 200 {object} orgdto.OrgAreaResp
// @Failure 400
// @Failure 500
// @Router /user/show/map [get]
func (u *UserCtrl) OrganizationInArea(w http.ResponseWriter, r *http.Request) {
	var req orgdto.OrgAreaReq
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
