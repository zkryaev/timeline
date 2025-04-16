package users

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/query"
	"timeline/internal/controller/scope"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/userdto"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type User interface {
	SearchOrgs(ctx context.Context, logger *zap.Logger, sreq *general.SearchReq) (*general.SearchResp, error)
	OrgsInArea(ctx context.Context, logger *zap.Logger, area *general.OrgAreaReq) (*general.OrgAreaResp, error)
	UserUpdate(ctx context.Context, logger *zap.Logger, user *userdto.UserUpdateReq) error
	User(ctx context.Context, logger *zap.Logger, id int) (*entity.User, error)
}

type UserCtrl struct {
	usecase    User
	Logger     *zap.Logger
	middleware middleware.Middleware
	settings   *scope.Settings
}

func New(usecase User, logger *zap.Logger, validator *validator.Validate, middleware middleware.Middleware, settings *scope.Settings) *UserCtrl {
	return &UserCtrl{
		usecase:    usecase,
		Logger:     logger,
		middleware: middleware,
		settings:   settings,
	}
}

// @Summary Get user
// @Description
// @Tags user
// @Accept  json
// @Produce  json
// @Param user_id query int true " "
// @Success 200 {object} entity.User
// @Failure 400
// @Failure 500
// @Router /users [get]
func (u *UserCtrl) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(u.settings, u.Logger, r.Context())
	var (
		userID = query.NewParamInt(scope.USER_ID, true)
	)
	params := query.NewParams(u.settings, userID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := u.usecase.User(r.Context(), logger, userID.Val)
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

// @Summary Update user
// @Description
// @Tags user
// @Accept  json
// @Produce  json
// @Param   req body userdto.UserUpdateReq true " "
// @Success 200 {object} userdto.UserUpdateReq
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /users [put]
func (u *UserCtrl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(u.settings, u.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(u.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req := &userdto.UserUpdateReq{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req.UserID = tdata.ID
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
// @Description Получить организацию удовлетворяющую параметрам поиска
// @Description `Если авторизация отключена: то *user_id* прокидывать в параметрах!`
// @Tags user
// @Accept  json
// @Produce  json
// @Param user_id query int true " "
// @Param limit query int true " "
// @Param page query int true " "
// @Param name query string false "Name of the organization to search for"
// @Param type query string false "Type of the organization"
// @Param sort_by query string false "Values: name/rate"
// @Success 200 {object} general.SearchResp
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /users/search/orgs [get]
func (u *UserCtrl) SearchOrganization(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(u.settings, u.Logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(u.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		userID  = query.NewParamInt(scope.USER_ID, true)
		limit   = query.NewParamInt(scope.LIMIT, true)
		page    = query.NewParamInt(scope.PAGE, true)
		orgName = query.NewParamString(scope.NAME, false)
		orgType = query.NewParamString(scope.TYPE, false)
		sortBy  = query.NewParamString(scope.SORT_BY, false)
	)
	params := query.NewParams(u.settings, limit, page, orgName, orgType, sortBy, userID)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	switch sortBy.Val {
	case scope.NAMESORT:
	case scope.RATESORT:
	default:
		sortBy.Val = ""
	}

	req := &general.SearchReq{
		Page:   page.Val,
		Limit:  limit.Val,
		Name:   orgName.Val,
		Type:   orgType.Val,
		SortBy: sortBy.Val,
		UserID: userID.Val,
	}
	if u.settings.EnableAuthorization {
		req.UserID = tdata.ID
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
// @Tags user
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
// @Router /users/orgmap [get]
func (u *UserCtrl) OrganizationInArea(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(u.settings, u.Logger, r.Context())
	var (
		minLat = query.NewParamFloat32(scope.MIN_LAT, true)
		minLon = query.NewParamFloat32(scope.MIN_LON, true)
		maxLat = query.NewParamFloat32(scope.MAX_LAT, true)
		maxLon = query.NewParamFloat32(scope.MAX_LON, true)
	)
	params := query.NewParams(u.settings, minLat, minLon, maxLat, maxLon)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	req := &general.OrgAreaReq{
		LeftLowerCorner: entity.Coordinates{
			Lat:  float64(minLat.Val),
			Long: float64(minLon.Val),
		},
		RightUpperCorner: entity.Coordinates{
			Lat:  float64(maxLat.Val),
			Long: float64(maxLon.Val),
		},
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
