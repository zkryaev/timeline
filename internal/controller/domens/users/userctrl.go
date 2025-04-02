package users

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/userdto"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type User interface {
	SearchOrgs(ctx context.Context, logger *zap.Logger, sreq *general.SearchReq) (*general.SearchResp, error)
	OrgsInArea(ctx context.Context, logger *zap.Logger, area *general.OrgAreaReq) (*general.OrgAreaResp, error)
	UserUpdate(ctx context.Context, logger *zap.Logger, user *userdto.UserUpdateReq) error
	User(ctx context.Context, logger *zap.Logger, id int) (*entity.User, error)
}

type UserCtrl struct {
	usecase User
	Logger  *zap.Logger
}

func New(usecase User, logger *zap.Logger, validator *validator.Validate) *UserCtrl {
	return &UserCtrl{
		usecase: usecase,
		Logger:  logger,
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
	uuid, _ := r.Context().Value("uuid").(string)
	logger := u.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "id")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := u.usecase.User(r.Context(), logger, params["id"])
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("User", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("User", zap.Error(err))
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

// @Summary Update
// @Description Update user information
// @Tags User
// @Accept  json
// @Produce  json
// @Param   request body userdto.UserUpdateReq true "New user info"
// @Success 200 {object} userdto.UserUpdateReq
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /users/update [put]
func (u *UserCtrl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := u.Logger.With(zap.String("uuid", uuid))
	req := &userdto.UserUpdateReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := u.usecase.UserUpdate(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("UserUpdate", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("UserUpdate", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
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
// @Param is_rate_sort query bool false "on/off rating sort"
// @Param is_name_sort query bool false "on/off name sort"
// @Success 200 {object} general.SearchResp
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /users/search/orgs [get]
func (u *UserCtrl) SearchOrganization(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := u.Logger.With(zap.String("uuid", uuid))
	query := map[string]bool{
		"limit":        true,
		"page":         true,
		"name":         false,
		"type":         false,
		"is_rate_sort": false,
		"is_name_sort": false,
	}
	if err := validation.IsQueryValid(r, query); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
	rateSort, _ := strconv.ParseBool(r.URL.Query().Get("is_rate_sort"))
	nameSort, _ := strconv.ParseBool(r.URL.Query().Get("is_name_sort"))

	req := &general.SearchReq{
		Page:       int(page),
		Limit:      int(limit),
		Name:       r.URL.Query().Get("name"),
		Type:       r.URL.Query().Get("type"),
		IsRateSort: rateSort,
		IsNameSort: nameSort,
	}
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := u.usecase.SearchOrgs(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("SearchOrgs", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("SearchOrgs", zap.Error(err))
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
// @Failure 404
// @Failure 500
// @Router /users/map/orgs [get]
func (u *UserCtrl) OrganizationInArea(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := u.Logger.With(zap.String("uuid", uuid))
	query := map[string]bool{
		"min_lat":  true,
		"min_long": true,
		"max_lat":  true,
		"max_long": true,
	}
	if err := validation.IsQueryValid(r, query); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
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
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := u.usecase.OrgsInArea(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("OrgsInArea", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("OrgsInArea", zap.Error(err))
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
