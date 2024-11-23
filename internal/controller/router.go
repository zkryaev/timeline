package controller

import (
	"timeline/internal/controller/auth"
	"timeline/internal/controller/domens"

	"github.com/gorilla/mux"
)

type Controllers struct {
	Auth *auth.AuthCtrl
	User *domens.UserCtrl
}

// Auth
const (
	authPrefix        = "/auth"
	authLogin         = "/login"
	authRegisterOrg   = "/register/org"
	authRegisterUser  = "/register/user"
	authRefreshToken  = "/refresh/token"
	authSendCodeRetry = "/send/code"
	authVerifyCode    = "/verify/code"
)

const (
	userPrefix     = "/user"
	userMapOrgs    = "/user/show/map"
	userSearchOrgs = "/user/find/orgs"
)

func InitRouter(controllersSet *Controllers) *mux.Router {
	r := mux.NewRouter()

	// Установка доменных контроллеров
	auth := controllersSet.Auth
	user := controllersSet.User

	r.Use(auth.Middleware.HandlerLogs)

	authRouter := r.NewRoute().PathPrefix(authPrefix).Subrouter()
	authRouter.HandleFunc(authLogin, auth.Login).Methods("POST")
	authRouter.HandleFunc(authRegisterOrg, auth.OrgRegister).Methods("POST")
	authRouter.HandleFunc(authRegisterUser, auth.UserRegister).Methods("POST")
	authRouter.HandleFunc(authRefreshToken, auth.UpdateAccessToken).Methods("PUT")
	authRouter.HandleFunc(authVerifyCode, auth.VerifyCode).Methods("POST")

	authProtectedRouter := r.NewRoute().PathPrefix("/auth").Subrouter()
	authProtectedRouter.Use(auth.Middleware.IsTokenValid)
	authProtectedRouter.HandleFunc(authSendCodeRetry, auth.SendCodeRetry).Methods("POST")

	userRouter := r.NewRoute().PathPrefix(userPrefix).Subrouter()
	userRouter.Use(auth.Middleware.IsTokenValid)
	userRouter.HandleFunc(userMapOrgs, user.OrganizationInArea).Methods("GET")
	userRouter.HandleFunc(userSearchOrgs, user.SearchOrganization).Methods("GET")
	return r
}
