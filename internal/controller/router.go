package controller

import (
	"timeline/internal/controller/auth"
	"timeline/internal/controller/domens/orgs"
	"timeline/internal/controller/domens/users"

	"github.com/gorilla/mux"
)

type Controllers struct {
	Auth *auth.AuthCtrl
	User *users.UserCtrl
	Org  *orgs.OrgCtrl
}

// General
const (
	health = "/health"
	v1     = "/v1"
)

// Auth
const (
	authPrefix        = "/auth"
	authLogin         = "/login"
	authRegisterOrg   = "/orgs"
	authRegisterUser  = "/users"
	authRefreshToken  = "/tokens/refresh"
	authSendCodeRetry = "/codes/send"
	authVerifyCode    = "/codes/verify"
)

// User
const (
	userPrefix     = "/users"
	userMapOrgs    = "/map/orgs"
	userSearchOrgs = "/search/orgs"
	userUpdate     = "/update"
	userGetInfo    = "/info/{id}"
)

const (
	orgPrefix          = "/orgs"
	orgGetInfo         = "/info/{id}"
	orgUpdate          = "/update"
	orgUpdateTimetable = "/{id}/timetable"
)

func InitRouter(controllersSet *Controllers) *mux.Router {
	r := mux.NewRouter()

	// Установка доменных контроллеров
	auth := controllersSet.Auth
	user := controllersSet.User
	org := controllersSet.Org

	r.Use(auth.Middleware.HandlerLogs)
	r.HandleFunc(health, HealthCheck)

	v1 := r.NewRoute().PathPrefix(v1).Subrouter()
	// Auth
	// !!!! Пока версия не продовая все ручки доступны без токенов !!!!
	authRouter := v1.NewRoute().PathPrefix(authPrefix).Subrouter()
	authRouter.HandleFunc(authLogin, auth.Login).Methods("POST")
	authRouter.HandleFunc(authRegisterOrg, auth.OrgRegister).Methods("POST")
	authRouter.HandleFunc(authRegisterUser, auth.UserRegister).Methods("POST")
	authRouter.HandleFunc(authRefreshToken, auth.UpdateAccessToken).Methods("PUT")
	authRouter.HandleFunc(authVerifyCode, auth.VerifyCode).Methods("POST")

	authProtectedRouter := v1.NewRoute().PathPrefix("/auth").Subrouter()
	// authProtectedRouter.Use(auth.Middleware.IsTokenValid)
	authProtectedRouter.HandleFunc(authSendCodeRetry, auth.SendCodeRetry).Methods("POST")

	// User
	userRouter := v1.NewRoute().PathPrefix(userPrefix).Subrouter()
	// userRouter.Use(auth.Middleware.IsTokenValid)
	userRouter.HandleFunc(userMapOrgs, user.OrganizationInArea).Methods("GET")
	userRouter.HandleFunc(userSearchOrgs, user.SearchOrganization).Methods("GET")
	userRouter.HandleFunc(userGetInfo, user.GetUserByID).Methods("GET")
	userRouter.HandleFunc(userUpdate, user.UpdateUser).Methods("PUT")
	// Org
	orgRouter := v1.NewRoute().PathPrefix(orgPrefix).Subrouter()
	// orgRouter.Use(auth.Middleware.IsTokenValid)
	orgRouter.HandleFunc(orgGetInfo, org.GetOrgByID).Methods("GET")
	orgRouter.HandleFunc(orgUpdate, org.UpdateOrg).Methods("PUT")
	orgRouter.HandleFunc(orgUpdateTimetable, org.UpdateOrgTimetable).Methods("PUT")

	return r
}
