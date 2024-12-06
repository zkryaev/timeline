package users

import (
	"context"
	"net/http"
	"strconv"
	"timeline/internal/controller/validation"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/userdto"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type User interface {
	SearchOrgs(ctx context.Context, sreq *general.SearchReq) (*general.SearchResp, error)
	OrgsInArea(ctx context.Context, area *general.OrgAreaReq) (*general.OrgAreaResp, error)
	UserUpdate(ctx context.Context, user *userdto.UserUpdateReq) error
	User(ctx context.Context, id int) (*entity.User, error)
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
// @Success 200 {object} userdto.UserUpdateReq
// @Failure 400
// @Failure 500
// @Router /users/update [put]
func (u *UserCtrl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	req := &userdto.UserUpdateReq{}
	if u.json.NewDecoder(r.Body).Decode(req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if err := u.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := u.usecase.UserUpdate(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Organization Searching
// @Description Get organizations that are satisfiered to search params
// @Tags User
// @Accept  json
// @Produce  json
// @Param limit query int true "Limit the number of results"
// @Param page query int true "Page number for pagination"
// @Param name query string false "Name of the organization to search for"
// @Param type query string false "Type of the organization"
// @Success 200 {object} general.SearchResp
// @Failure 400
// @Failure 500
// @Router /users/search/orgs [get]
func (u *UserCtrl) SearchOrganization(w http.ResponseWriter, r *http.Request) {
	query := map[string]bool{
		"limit": true,
		"page":  true,
		"name":  false,
		"type":  false,
	}
	if !validation.IsQueryValid(r, query) {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)

	req := &general.SearchReq{
		Page:  int(page),
		Limit: int(limit),
		Name:  r.URL.Query().Get("name"),
		Type:  r.URL.Query().Get("type"),
	}
	if err := u.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, err := u.usecase.SearchOrgs(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if u.json.NewEncoder(w).Encode(data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
}

// @Summary Show Organizations on Map
// @Description Get organizations that located at the given area
// @Tags User
// @Accept  json
// @Produce  json
// @Param min_lat query float32 true "Minimum latitude for the search area"
// @Param min_long query float32 true "Minimum longitude for the search area"
// @Param max_lat query float32 true "Maximum latitude for the search area"
// @Param max_long query float32 true "Maximum longitude for the search area"
// @Success 200 {object} general.OrgAreaResp
// @Failure 400
// @Failure 500
// @Router /users/map/orgs [get]
func (u *UserCtrl) OrganizationInArea(w http.ResponseWriter, r *http.Request) {
	query := map[string]bool{
		"min_lat":  true,
		"min_long": true,
		"max_lat":  true,
		"max_long": true,
	}
	if !validation.IsQueryValid(r, query) {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}
	minLat, _ := strconv.ParseFloat(r.URL.Query().Get("min_lat"), 64)
	minLong, _ := strconv.ParseFloat(r.URL.Query().Get("min_long"), 64)
	maxLat, _ := strconv.ParseFloat(r.URL.Query().Get("max_lat"), 64)
	maxLong, _ := strconv.ParseFloat(r.URL.Query().Get("max_long"), 64)
	req := &general.OrgAreaReq{
		LeftLowerCorner: entity.Coordinates{
			Lat:  minLat,
			Long: minLong,
		},
		RightUpperCorner: entity.Coordinates{
			Lat:  maxLat,
			Long: maxLong,
		},
	}
	// валидация полей
	if err := u.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := u.usecase.OrgsInArea(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// отдаем токен
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if u.json.NewEncoder(w).Encode(data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
}

// @Summary Get user
// @Description Get user by his id
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "user_id"
// @Success 200 {object} entity.User
// @Failure 400
// @Failure 500
// @Router /users/info/{id} [get]
func (u *UserCtrl) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := u.usecase.User(r.Context(), params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if u.json.NewEncoder(w).Encode(data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}

}
