package controller

import (
	"timeline/internal/controller/auth"

	"github.com/gorilla/mux"
)

type Controllers struct {
	Auth *auth.AuthCtrl
}

const (
	authPrefix        = "/auth"
	authLogin         = "/login"
	authRegisterOrg   = "/register/org"
	authRegisterUser  = "/register/user"
	authRefreshToken  = "/refresh/token"
	authSendCodeRetry = "/send/code"
	authVerifyCode    = "/verify/code"
)

func InitRouter(controllersSet *Controllers) *mux.Router {
	r := mux.NewRouter()

	// Установка доменных контроллеров
	auth := controllersSet.Auth

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
	return r
}
