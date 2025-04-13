package controller

import (
	"timeline/internal/controller/auth"
	"timeline/internal/controller/domens/orgs"
	"timeline/internal/controller/domens/records"
	"timeline/internal/controller/domens/users"
	"timeline/internal/controller/s3"
	"timeline/internal/controller/settings"

	"github.com/gorilla/mux"
)

type Controllers struct {
	Auth   *auth.AuthCtrl
	User   *users.UserCtrl
	Org    *orgs.OrgCtrl
	Record *records.RecordCtrl
	S3     *s3.S3Ctrl
}

func InitRouter(controllersSet *Controllers, routes settings.Routes) *mux.Router {
	r := mux.NewRouter()

	// Установка доменных контроллеров
	auth := controllersSet.Auth
	user := controllersSet.User
	org := controllersSet.Org
	rec := controllersSet.Record
	s3 := controllersSet.S3

	// TODO: PROD s := r.Host("www.example.com").Subrouter()
	r.Use(auth.Middleware.HandlerLogs)
	r.HandleFunc(settings.PathHealth, HealthCheck)
	V1 := r.NewRoute().PathPrefix(settings.V1).Subrouter()
	// Auth
	authmux := V1.NewRoute().PathPrefix(settings.PathAuth).Subrouter()
	authmux.HandleFunc(settings.PathLogin, auth.Login).Methods(routes[settings.PathLogin].Methods.Get(settings.POST)...)
	authmux.HandleFunc(settings.PathRegistration, auth.OrgRegister).Methods(routes[settings.PathRegistration].Methods.Get(settings.POST)...)
	authmux.HandleFunc(settings.PathToken, auth.UpdateAccessToken).Methods(routes[settings.PathToken].Methods.Get(settings.PUT)...)

	Protected := V1.NewRoute().Subrouter()
	Protected.Use(auth.Middleware.Authorization)

	authmuxProtected := Protected.NewRoute().PathPrefix(settings.PathAuth).Subrouter()
	authmuxProtected.HandleFunc(settings.PathCode, auth.SendCodeRetry).Methods(routes[settings.PathCode].Methods.Get(settings.POST)...)
	authmuxProtected.HandleFunc(settings.PathCode, auth.VerifyCode).Methods(routes[settings.PathCode].Methods.Get(settings.PUT)...)

	// users
	Protected.HandleFunc(settings.PathUsers, user.GetUserByID).Methods(routes[settings.PathUsers].Methods.Get(settings.GET)...)
	Protected.HandleFunc(settings.PathUsers, user.UpdateUser).Methods(routes[settings.PathUsers].Methods.Get(settings.PUT)...)
	// orgs
	Protected.HandleFunc(settings.PathOrgs, org.GetOrgByID).Methods(routes[settings.PathOrgs].Methods.Get(settings.GET)...)
	Protected.HandleFunc(settings.PathOrgs, org.UpdateOrg).Methods(routes[settings.PathOrgs].Methods.Get(settings.PUT)...)

	usermuxProtected := Protected.NewRoute().PathPrefix(settings.PathUsers).Subrouter()
	// users/orgmap
	usermuxProtected.HandleFunc(settings.PathMapOrgs, user.OrganizationInArea).Methods(routes[settings.PathMapOrgs].Methods.Get(settings.GET)...)
	// users/search/org
	usermuxProtected.HandleFunc(settings.PathSearchOrgs, user.SearchOrganization).Methods(routes[settings.PathSearchOrgs].Methods.Get(settings.GET)...)

	orgmuxProtected := Protected.NewRoute().PathPrefix(settings.PathOrgs).Subrouter()
	// orgs/timetables
	orgmuxProtected.HandleFunc(settings.PathTimetables, org.TimetableAdd).Methods(routes[settings.PathTimetables].Methods.Get(settings.POST)...)
	orgmuxProtected.HandleFunc(settings.PathTimetables, org.Timetable).Methods(routes[settings.PathTimetables].Methods.Get(settings.GET)...)
	orgmuxProtected.HandleFunc(settings.PathTimetables, org.TimetableUpdate).Methods(routes[settings.PathTimetables].Methods.Get(settings.PUT)...)
	orgmuxProtected.HandleFunc(settings.PathTimetables, org.TimetableDelete).Methods(routes[settings.PathTimetables].Methods.Get(settings.DELETE)...)
	// orgs/services
	orgmuxProtected.HandleFunc(settings.PathServices, org.ServiceAdd).Methods(routes[settings.PathServices].Methods.Get(settings.POST)...)
	orgmuxProtected.HandleFunc(settings.PathServices, org.Service).Methods(routes[settings.PathServices].Methods.Get(settings.GET)...)
	orgmuxProtected.HandleFunc(settings.PathServices, org.ServiceUpdate).Methods(routes[settings.PathServices].Methods.Get(settings.PUT)...)
	orgmuxProtected.HandleFunc(settings.PathServices, org.ServiceDelete).Methods(routes[settings.PathServices].Methods.Get(settings.DELETE)...)
	// orgs/workers
	orgmuxProtected.HandleFunc(settings.PathWorkers, org.WorkerAdd).Methods(routes[settings.PathWorkers].Methods.Get(settings.POST)...)
	orgmuxProtected.HandleFunc(settings.PathWorkers, org.Worker).Methods(routes[settings.PathWorkers].Methods.Get(settings.GET)...)
	orgmuxProtected.HandleFunc(settings.PathWorkers, org.WorkerUpdate).Methods(routes[settings.PathWorkers].Methods.Get(settings.PUT)...)
	orgmuxProtected.HandleFunc(settings.PathWorkers, org.WorkerDelete).Methods(routes[settings.PathWorkers].Methods.Get(settings.DELETE)...)
	// orgs/workers/slots
	orgmuxProtected.HandleFunc(settings.PathWorkersSlots, org.Slots).Methods(routes[settings.PathWorkersSlots].Methods.Get(settings.GET)...)
	orgmuxProtected.HandleFunc(settings.PathWorkersSlots, org.UpdateSlot).Methods(routes[settings.PathWorkersSlots].Methods.Get(settings.PUT)...)
	// orgs/workers/services
	orgmuxProtected.HandleFunc(settings.PathWorkersServices, org.WorkerAssignService).Methods(routes[settings.PathWorkersServices].Methods.Get(settings.POST)...)
	orgmuxProtected.HandleFunc(settings.PathWorkersServices, org.WorkerUnAssignService).Methods(routes[settings.PathWorkersServices].Methods.Get(settings.DELETE)...)
	orgmuxProtected.HandleFunc(settings.PathWorkersServices, org.ServiceWorkerList).Methods(routes[settings.PathWorkersServices].Methods.Get(settings.GET)...)
	// orgs/workers/schedules
	orgmuxProtected.HandleFunc(settings.PathWorkersSchedules, org.AddWorkerSchedule).Methods(routes[settings.PathWorkersSchedules].Methods.Get(settings.POST)...)
	orgmuxProtected.HandleFunc(settings.PathWorkersSchedules, org.WorkerSchedule).Methods(routes[settings.PathWorkersSchedules].Methods.Get(settings.GET)...)
	orgmuxProtected.HandleFunc(settings.PathWorkersSchedules, org.UpdateWorkerSchedule).Methods(routes[settings.PathWorkersSchedules].Methods.Get(settings.PUT)...)
	orgmuxProtected.HandleFunc(settings.PathWorkersSchedules, org.DeleteWorkerSchedule).Methods(routes[settings.PathWorkersSchedules].Methods.Get(settings.DELETE)...)

	// records
	Protected.HandleFunc(settings.PathRecords, rec.RecordAdd).Methods(routes[settings.PathRecords].Methods.Get(settings.POST)...)
	Protected.HandleFunc(settings.PathRecords, rec.Record).Methods(routes[settings.PathRecords].Methods.Get(settings.GET)...)
	Protected.HandleFunc(settings.PathRecords, rec.RecordCancel).Methods(routes[settings.PathRecords].Methods.Get(settings.PUT)...)

	recmuxProtected := Protected.NewRoute().PathPrefix(settings.PathRecords).Subrouter()
	// records/feedbacks
	recmuxProtected.HandleFunc(settings.PathFeedback, rec.FeedbackSet).Methods(routes[settings.PathFeedback].Methods.Get(settings.POST)...)
	recmuxProtected.HandleFunc(settings.PathFeedback, rec.Feedbacks).Methods(routes[settings.PathFeedback].Methods.Get(settings.GET)...)
	recmuxProtected.HandleFunc(settings.PathFeedback, rec.FeedbackUpdate).Methods(routes[settings.PathFeedback].Methods.Get(settings.PUT)...)
	recmuxProtected.HandleFunc(settings.PathFeedback, rec.FeedbackDelete).Methods(routes[settings.PathFeedback].Methods.Get(settings.DELETE)...)

	// media
	Protected.HandleFunc(settings.PathMedia, s3.Upload).Methods(routes[settings.PathFeedback].Methods.Get(settings.POST)...)
	Protected.HandleFunc(settings.PathMedia, s3.Download).Methods(routes[settings.PathFeedback].Methods.Get(settings.GET)...)
	Protected.HandleFunc(settings.PathMedia, s3.Delete).Methods(routes[settings.PathFeedback].Methods.Get(settings.DELETE)...)
	return r
}
