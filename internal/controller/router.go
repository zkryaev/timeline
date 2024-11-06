package controller

import (
	"net/http"
	"timeline/internal/controller/auth"

	"github.com/gorilla/mux"
)

type Controllers struct {
	Auth *auth.AuthCtrl
}

const (
	authPrefix       = "/auth"
	authLogin        = "/login"
	authRegisterOrg  = "/register/org"
	authRegisterUser = "/register/user"
	authRefreshToken = "/token/update"
)

func InitRouter(controllersSet *Controllers) *mux.Router {
	r := mux.NewRouter()

	auth := controllersSet.Auth

	authRouter := r.NewRoute().PathPrefix(authPrefix).Subrouter()
	authRouter.HandleFunc(authLogin, auth.Login).Methods("POST")
	authRouter.HandleFunc(authRegisterOrg, auth.OrgRegister).Methods("POST")
	authRouter.HandleFunc(authRegisterUser, auth.UserRegister).Methods("POST")
	authRouter.HandleFunc(authRefreshToken, auth.Middleware.IsRefreshToken(http.HandlerFunc(auth.UpdateAccessToken)).ServeHTTP).Methods("PUT")
	return r
}
