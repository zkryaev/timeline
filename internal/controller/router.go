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

// User
const (
	userPrefix     = "/user"
	userMapOrgs    = "/show/map"
	userSearchOrgs = "/find/orgs"
	userUpdate     = "/update"
)

const (
	orgPrefix = "/org"
	orgUpdate = "/update"
)

func InitRouter(controllersSet *Controllers) *mux.Router {
	r := mux.NewRouter()

	// Установка доменных контроллеров
	auth := controllersSet.Auth
	user := controllersSet.User
	org := controllersSet.Org

	r.Use(auth.Middleware.HandlerLogs)

	// Auth
	// !!!! Пока версия не продовая все ручки доступны без токенов !!!!
	authRouter := r.NewRoute().PathPrefix(authPrefix).Subrouter()
	authRouter.HandleFunc(authLogin, auth.Login).Methods("POST")
	authRouter.HandleFunc(authRegisterOrg, auth.OrgRegister).Methods("POST")
	authRouter.HandleFunc(authRegisterUser, auth.UserRegister).Methods("POST")
	authRouter.HandleFunc(authRefreshToken, auth.UpdateAccessToken).Methods("PUT")
	authRouter.HandleFunc(authVerifyCode, auth.VerifyCode).Methods("POST")

	authProtectedRouter := r.NewRoute().PathPrefix("/auth").Subrouter()
	authProtectedRouter.Use(auth.Middleware.IsTokenValid)
	authProtectedRouter.HandleFunc(authSendCodeRetry, auth.SendCodeRetry).Methods("POST")

	// User
	userRouter := r.NewRoute().PathPrefix(userPrefix).Subrouter()
	// userRouter.Use(auth.Middleware.IsTokenValid)
	userRouter.HandleFunc(userMapOrgs, user.OrganizationInArea).Methods("GET")
	userRouter.HandleFunc(userSearchOrgs, user.SearchOrganization).Methods("GET")
	userRouter.HandleFunc(userUpdate, user.UpdateUser).Methods("PUT")
	// Org
	orgRouter := r.NewRoute().PathPrefix(orgPrefix).Subrouter()
	// orgRouter.Use(auth.Middleware.IsTokenValid)
	orgRouter.HandleFunc(orgUpdate, org.UpdateOrg).Methods("PUT")

	return r
}
