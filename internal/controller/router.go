package controller

import (
	"timeline/internal/controller/auth"
	"timeline/internal/controller/domens/orgs"
	"timeline/internal/controller/domens/records"
	"timeline/internal/controller/domens/users"
	"timeline/internal/controller/s3"

	"github.com/gorilla/mux"
)

type Controllers struct {
	Auth   *auth.AuthCtrl
	User   *users.UserCtrl
	Org    *orgs.OrgCtrl
	Record *records.RecordCtrl
	S3     *s3.S3Ctrl
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

// Org
const (
	orgPrefix  = "/orgs"
	orgGetInfo = "/info/{id}"
	orgUpdate  = "/update"
	// Timetables
	timetable   = "/timetable"
	timetableID = "/{orgID}/timetable"
	// Workers
	worker         = "/workers"
	workerID       = "/{orgID}/workers/{workerID}"
	workerList     = "/{orgID}/workers"
	workerAssign   = "/workers/service"
	workerUnAssign = "/{orgID}/workers/service/{workerID}/{serviceID}"
	// Services
	service        = "/services"
	serviceID      = "/{orgID}/services/{serviceID}"
	serviceWorkers = "/{orgID}/services/{serviceID}/workers"
	serviceList    = "/{orgID}/services"
	// Schedule
	schedule        = "/schedules"
	scheduleWorkers = "/{orgID}/schedules"
	scheduleDelete  = "/{orgID}/schedules/{workerID}"

	// Slots
	slots       = "/{orgID}/slots"
	slotsWorker = "/{orgID}/slots/workers/{workerID}"
)

// Record
const (
	record     = "/records"
	recordAdd  = "/creation"
	recordID   = "/info/{recordID}"
	recordList = "/list"
	// Feedback
	feedback   = "/feedbacks"
	feedbackID = "/feedbacks/info"
)

// S3
const (
	media = "/media"
)

func InitRouter(controllersSet *Controllers) *mux.Router {
	r := mux.NewRouter()

	// Установка доменных контроллеров
	auth := controllersSet.Auth
	user := controllersSet.User
	org := controllersSet.Org
	rec := controllersSet.Record
	s3 := controllersSet.S3

	r.Use(auth.Middleware.HandlerLogs)
	r.HandleFunc(health, HealthCheck)

	v1 := r.NewRoute().PathPrefix(v1).Subrouter()
	// !!!! Пока версия не продовая все ручки доступны без токенов !!!!
	// Auth
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
	// Timetable
	orgRouter.HandleFunc(timetable, org.TimetableAdd).Methods("POST")
	orgRouter.HandleFunc(timetable, org.TimetableUpdate).Methods("PUT")
	orgRouter.HandleFunc(timetableID, org.Timetable).Methods("GET")
	orgRouter.HandleFunc(timetableID, org.TimetableDelete).Methods("DELETE")

	// Workers
	orgRouter.HandleFunc(worker, org.WorkerAdd).Methods("POST")
	orgRouter.HandleFunc(worker, org.WorkerUpdate).Methods("PUT")
	orgRouter.HandleFunc(workerID, org.WorkerDelete).Methods("DELETE")
	orgRouter.HandleFunc(workerID, org.Worker).Methods("GET")
	orgRouter.HandleFunc(workerList, org.WorkerList).Methods("GET")
	orgRouter.HandleFunc(workerAssign, org.WorkerAssignService).Methods("POST")
	orgRouter.HandleFunc(workerUnAssign, org.WorkerUnAssignService).Methods("DELETE")
	// Services
	orgRouter.HandleFunc(service, org.ServiceAdd).Methods("POST")
	orgRouter.HandleFunc(service, org.ServiceUpdate).Methods("PUT")
	orgRouter.HandleFunc(serviceID, org.Service).Methods("GET")
	orgRouter.HandleFunc(serviceID, org.ServiceDelete).Methods("DELETE")
	orgRouter.HandleFunc(serviceWorkers, org.ServiceWorkerList).Methods("GET")
	orgRouter.HandleFunc(serviceList, org.ServiceList).Methods("GET")
	// Schedule
	orgRouter.HandleFunc(schedule, org.AddWorkerSchedule).Methods("POST")
	orgRouter.HandleFunc(schedule, org.UpdateWorkerSchedule).Methods("PUT")
	orgRouter.HandleFunc(scheduleWorkers, org.WorkerSchedule).Methods("GET")
	orgRouter.HandleFunc(scheduleDelete, org.DeleteWorkerSchedule).Methods("DELETE")
	// Slots
	orgRouter.HandleFunc(slotsWorker, org.Slots).Methods("GET")
	orgRouter.HandleFunc(slots, org.UpdateSlot).Methods("PUT")

	// Records
	recRouter := v1.NewRoute().PathPrefix(record).Subrouter()
	recRouter.HandleFunc(recordAdd, rec.RecordAdd).Methods("POST")
	recRouter.HandleFunc(recordID, rec.Record).Methods("GET")
	recRouter.HandleFunc(recordList, rec.RecordList).Methods("GET")
	recRouter.HandleFunc(recordID, rec.RecordDelete).Methods("DELETE")
	// Feedbacks
	recRouter.HandleFunc(feedback, rec.FeedbackSet).Methods("POST")
	recRouter.HandleFunc(feedback, rec.FeedbackUpdate).Methods("PUT")
	recRouter.HandleFunc(feedbackID, rec.Feedbacks).Methods("GET")
	recRouter.HandleFunc(feedbackID, rec.FeedbackDelete).Methods("DELETE")

	s3Router := v1.NewRoute().Subrouter()
	s3Router.HandleFunc(media, s3.Upload).Methods("POST")
	s3Router.HandleFunc(media, s3.Download).Methods("GET")
	s3Router.HandleFunc(media, s3.Delete).Methods("DELETE")
	return r
}
