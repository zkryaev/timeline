package controller

import (
	"fmt"
	"net/http"
	"timeline/internal/controller/auth"
	"timeline/internal/controller/domens/orgs"
	"timeline/internal/controller/domens/records"
	"timeline/internal/controller/domens/users"
	"timeline/internal/controller/s3"
	"timeline/internal/controller/scope"

	"github.com/gorilla/mux"
)

type Controllers struct {
	Auth   *auth.AuthCtrl
	User   *users.UserCtrl
	Org    *orgs.OrgCtrl
	Record *records.RecordCtrl
	S3     *s3.S3Ctrl
}

func InitRouter(controllersSet *Controllers, routes scope.Routes, settings *scope.Settings) *mux.Router {
	r := mux.NewRouter()

	// Установка доменных контроллеров
	auth := controllersSet.Auth
	user := controllersSet.User
	org := controllersSet.Org
	rec := controllersSet.Record
	s3 := controllersSet.S3

	// s := r.Host("www.timeline.com").Subrouter() // TODO future

	r.Use(auth.Middleware.HandlerLogs)

	r.HandleFunc(scope.PathHealth, HealthCheck).Methods(http.MethodGet)

	v1 := r.PathPrefix(scope.ApiVersionV1).Subrouter()

	// v1/auth
	authmux := v1.PathPrefix(scope.PathAuth).Subrouter()
	authmux.HandleFunc(scope.PathLogin, auth.Login).Methods(routes[scope.PathLogin].Methods.Get(scope.POST)...)
	authmux.HandleFunc(scope.PathUsersRegistration, auth.UserRegister).Methods(routes[scope.PathUsersRegistration].Methods.Get(scope.POST)...)
	authmux.HandleFunc(scope.PathOrgsRegistration, auth.OrganizationRegister).Methods(routes[scope.PathOrgsRegistration].Methods.Get(scope.POST)...)
	authmux.HandleFunc(scope.PathToken, auth.PutAccessToken).Methods(routes[scope.PathToken].Methods.Get(scope.PUT)...)

	Protected := v1.NewRoute().Subrouter()
	if settings.EnableAuthorization {
		Protected.Use(auth.Middleware.Authorization)
	}

	// v1/auth/codes
	authmuxProtected := Protected.PathPrefix(scope.PathAuth).Subrouter()
	authmuxProtected.HandleFunc(scope.PathCode, auth.CodeSend).Methods(routes[scope.PathCode].Methods.Get(scope.POST)...)
	authmuxProtected.HandleFunc(scope.PathCode, auth.CodeConfirm).Methods(routes[scope.PathCode].Methods.Get(scope.PUT)...)

	// v1/users
	usermuxProtected := Protected.PathPrefix(scope.PathUsers).Subrouter()
	usermuxProtected.HandleFunc("", user.GetUser).Methods(routes[scope.PathUsers].Methods.Get(scope.GET)...)
	usermuxProtected.HandleFunc("", user.UpdateUser).Methods(routes[scope.PathUsers].Methods.Get(scope.PUT)...)

	// v1/users/orgmap
	usermuxProtected.HandleFunc(scope.PathMapOrgs, user.OrganizationInArea).Methods(routes[scope.PathMapOrgs].Methods.Get(scope.GET)...)

	// v1/users/search/org
	usermuxProtected.HandleFunc(scope.PathSearchOrgs, user.SearchOrganization).Methods(routes[scope.PathSearchOrgs].Methods.Get(scope.GET)...)

	// orgs
	orgmuxProtected := Protected.PathPrefix(scope.PathOrgs).Subrouter()
	orgmuxProtected.HandleFunc("", org.GetOrganization).Methods(routes[scope.PathOrgs].Methods.Get(scope.GET)...)
	orgmuxProtected.HandleFunc("", org.PutOrganization).Methods(routes[scope.PathOrgs].Methods.Get(scope.PUT)...)

	// orgs/timetables
	orgmuxProtected.HandleFunc(scope.PathTimetables, org.TimetableAdd).Methods(routes[scope.PathTimetables].Methods.Get(scope.POST)...)
	orgmuxProtected.HandleFunc(scope.PathTimetables, org.Timetable).Methods(routes[scope.PathTimetables].Methods.Get(scope.GET)...)
	orgmuxProtected.HandleFunc(scope.PathTimetables, org.TimetableUpdate).Methods(routes[scope.PathTimetables].Methods.Get(scope.PUT)...)
	orgmuxProtected.HandleFunc(scope.PathTimetables, org.TimetableDelete).Methods(routes[scope.PathTimetables].Methods.Get(scope.DELETE)...)

	// orgs/services
	orgmuxProtected.HandleFunc(scope.PathServices, org.ServiceAdd).Methods(routes[scope.PathServices].Methods.Get(scope.POST)...)
	orgmuxProtected.HandleFunc(scope.PathServices, org.Service).Methods(routes[scope.PathServices].Methods.Get(scope.GET)...)
	orgmuxProtected.HandleFunc(scope.PathServices, org.ServiceUpdate).Methods(routes[scope.PathServices].Methods.Get(scope.PUT)...)
	orgmuxProtected.HandleFunc(scope.PathServices, org.ServiceDelete).Methods(routes[scope.PathServices].Methods.Get(scope.DELETE)...)

	// orgs/workers
	orgmuxProtected.HandleFunc(scope.PathWorkers, org.WorkerAdd).Methods(routes[scope.PathWorkers].Methods.Get(scope.POST)...)
	orgmuxProtected.HandleFunc(scope.PathWorkers, org.Workers).Methods(routes[scope.PathWorkers].Methods.Get(scope.GET)...)
	orgmuxProtected.HandleFunc(scope.PathWorkers, org.WorkerUpdate).Methods(routes[scope.PathWorkers].Methods.Get(scope.PUT)...)
	orgmuxProtected.HandleFunc(scope.PathWorkers, org.WorkerDelete).Methods(routes[scope.PathWorkers].Methods.Get(scope.DELETE)...)

	// orgs/workers/slots
	orgmuxProtected.HandleFunc(scope.PathWorkersSlots, org.Slots).Methods(routes[scope.PathWorkersSlots].Methods.Get(scope.GET)...)

	// orgs/workers/services
	orgmuxProtected.HandleFunc(scope.PathWorkersServices, org.WorkerAssignService).Methods(routes[scope.PathWorkersServices].Methods.Get(scope.POST)...)
	orgmuxProtected.HandleFunc(scope.PathWorkersServices, org.WorkersServices).Methods(routes[scope.PathWorkersServices].Methods.Get(scope.GET)...)
	orgmuxProtected.HandleFunc(scope.PathWorkersServices, org.WorkerUnassignService).Methods(routes[scope.PathWorkersServices].Methods.Get(scope.DELETE)...)

	// orgs/workers/schedules
	orgmuxProtected.HandleFunc(scope.PathWorkersSchedules, org.AddWorkerSchedule).Methods(routes[scope.PathWorkersSchedules].Methods.Get(scope.POST)...)
	orgmuxProtected.HandleFunc(scope.PathWorkersSchedules, org.WorkersSchedule).Methods(routes[scope.PathWorkersSchedules].Methods.Get(scope.GET)...)
	orgmuxProtected.HandleFunc(scope.PathWorkersSchedules, org.UpdateWorkerSchedule).Methods(routes[scope.PathWorkersSchedules].Methods.Get(scope.PUT)...)
	orgmuxProtected.HandleFunc(scope.PathWorkersSchedules, org.DeleteWorkerSchedule).Methods(routes[scope.PathWorkersSchedules].Methods.Get(scope.DELETE)...)

	// records
	recmuxProtected := Protected.PathPrefix(scope.PathRecords).Subrouter()
	recmuxProtected.HandleFunc("", rec.RecordAdd).Methods(routes[scope.PathRecords].Methods.Get(scope.POST)...)
	recmuxProtected.HandleFunc("", rec.Record).Methods(routes[scope.PathRecords].Methods.Get(scope.GET)...)
	recmuxProtected.HandleFunc("", rec.RecordCancel).Methods(routes[scope.PathRecords].Methods.Get(scope.PUT)...)

	// records/feedbacks
	recmuxProtected.HandleFunc(scope.PathFeedback, rec.FeedbackSet).Methods(routes[scope.PathFeedback].Methods.Get(scope.POST)...)
	recmuxProtected.HandleFunc(scope.PathFeedback, rec.Feedbacks).Methods(routes[scope.PathFeedback].Methods.Get(scope.GET)...)
	recmuxProtected.HandleFunc(scope.PathFeedback, rec.FeedbackUpdate).Methods(routes[scope.PathFeedback].Methods.Get(scope.PUT)...)
	recmuxProtected.HandleFunc(scope.PathFeedback, rec.FeedbackDelete).Methods(routes[scope.PathFeedback].Methods.Get(scope.DELETE)...)

	// media
	if settings.EnableRepoS3 {
		Protected.HandleFunc(scope.PathMedia, s3.Upload).Methods(routes[scope.PathMedia].Methods.Get(scope.POST)...)
		Protected.HandleFunc(scope.PathMedia, s3.Download).Methods(routes[scope.PathMedia].Methods.Get(scope.GET)...)
		Protected.HandleFunc(scope.PathMedia, s3.Delete).Methods(routes[scope.PathMedia].Methods.Get(scope.DELETE)...)
	}

	fmt.Println("=== Router name:", "r")
	printAllRoutes(r)

	return r
}

func printAllRoutes(router *mux.Router) {
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			// Если нет методов (например, подроутер), просто печатаем путь
			fmt.Printf("Path: %-40s (subrouter)\n", pathTemplate)
			return nil
		}

		fmt.Printf("Path: %-40s Methods: %v\n", pathTemplate, methods)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking routes: %v\n", err)
	}
}
