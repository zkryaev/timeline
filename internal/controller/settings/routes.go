package settings

import (
	"net/http"
	"timeline/internal/controller/perms"
)

const (
	V1 = "/v1"
)

const (
	// Utility endpoint
	PathHealth = "/health" // docker health handler

	// Auth endpoint
	PathAuth         = "/auth"
	PathLogin        = "/login"
	PathRegistration = "/registration"
	PathToken        = "/token"
	PathCode         = "/code"

	// User endpoint
	PathUsers      = "/users"
	PathMapOrgs    = PathUsers + "/orgmap"
	PathSearchOrgs = PathUsers + "/search/org"

	// Organization endpoint
	PathOrgs             = "/orgs"
	PathWorkers          = "/workers"
	PathServices         = "/services"
	PathTimetables       = "/timetables"
	PathWorkersSlots     = PathWorkers + "/slots"
	PathWorkersServices  = PathWorkers + "/services"
	PathWorkersSchedules = PathWorkers + "/schedules"

	// Record endpoint
	PathRecords  = "/records"
	PathFeedback = PathRecords + "/feedbacks"

	// S3 endpoint
	PathMedia = "/media"
)

var PathList = []string{
	// auth
	PathAuth,
	PathLogin,
	PathRegistration,
	PathToken,
	PathCode,
	// users
	PathUsers,
	PathMapOrgs,
	PathSearchOrgs,
	// org
	PathOrgs,
	PathServices,
	PathTimetables,
	// org/worker/...
	PathWorkers,
	PathWorkersSlots,
	PathWorkersServices,
	PathWorkersSchedules,
	// records/...
	PathRecords,
	PathFeedback,
	// media
	PathMedia,
}

type MethodList []string

func newMethodsMap(s *Settings, methods ...string) MethodList {
	m := make(MethodList, len(methods))
	for i := range m {
		if ind, ok := s.SupportedMethodsMap[methods[i]]; ok {
			m[ind] = methods[i]
		}
	}
	return m
}

// Return methods specified by method enum in settings.go
//
// Call endpoint.Methods... - to get all methods
func (mp MethodList) Get(enum ...uint8) []string {
	list := make([]string, len(enum))
	for i := range enum {
		list = append(list, mp[enum[i]])
	}
	return list
}

type endpoint struct {
	Path    string
	Methods MethodList
	perms   perms.PermissionBits
}

func NewEndpointFromPath(s *Settings, path string) endpoint {
	mdata := endpoint{}
	switch path {
	// /auth  [ Everybody: All ]
	// /auth/login
	case PathLogin:
		mdata.Methods = newMethodsMap(s, http.MethodPost)
		mdata.perms = perms.GrantPermissions(perms.CREATE, perms.CREATE)
	// /auth/registration
	case PathRegistration:
		mdata.Methods = newMethodsMap(s, http.MethodPost)
		mdata.perms = perms.GrantPermissions(perms.CREATE, perms.CREATE)
	// /auth/token
	case PathToken:
		mdata.Methods = newMethodsMap(s, http.MethodPost)
		mdata.perms = perms.GrantPermissions(perms.CREATE+perms.UPDATE, perms.CREATE+perms.UPDATE)
	// /auth/code
	case PathCode:
		mdata.Methods = newMethodsMap(s, http.MethodPost, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.CREATE, perms.CREATE)
	// /users
	case PathUsers:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.READ+perms.UPDATE, perms.NONE)
	// /users/orgmap
	case PathMapOrgs:
		mdata.Methods = newMethodsMap(s, http.MethodGet)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.NONE)
	// /users/search/org
	case PathSearchOrgs:
		mdata.Methods = newMethodsMap(s, http.MethodGet)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.NONE)
	// /orgs  [ Users: Only Read | Orgs: All ]
	case PathOrgs:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.NONE, perms.READ+perms.UPDATE)
	// /orgs/timetables
	case PathTimetables:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.ALL)
	// /orgs/services
	case PathServices:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.ALL)
	// /orgs/workers
	case PathWorkers:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.ALL)
	// /orgs/workers/services
	case PathWorkersServices:
		mdata.Methods = newMethodsMap(s, http.MethodPost, http.MethodDelete)
		mdata.perms = perms.GrantPermissions(perms.NONE, perms.CREATE+perms.DELETE)
	// /orgs/workers/schedules
	case PathWorkersSchedules:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.ALL)
	// /orgs/workers/slots
	case PathWorkersSlots:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.READ+perms.UPDATE)
	// /records  [ Users: All | Orgs: Only Read ]
	case PathRecords:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPost, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.CREATE+perms.READ+perms.UPDATE, perms.READ+perms.UPDATE)
	// /records/feedbacks
	case PathFeedback:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		mdata.perms = perms.GrantPermissions(perms.ALL, perms.READ)
	// /media  [ Everybody: All ]
	case PathMedia:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPost, http.MethodDelete)
		mdata.perms = perms.GrantPermissions(perms.CREATE+perms.READ+perms.DELETE, perms.CREATE+perms.READ+perms.DELETE)
	}
	return mdata
}

type Routes map[string]endpoint

func NewDefaultRoutes(settings *Settings) Routes {
	r := make(Routes, len(PathList))
	for _, path := range PathList {
		r[path] = NewEndpointFromPath(settings, path)
	}
	return r
}
