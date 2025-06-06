package scope

import (
	"fmt"
	"net/http"
	"strings"
	"timeline/internal/controller/perms"
	"timeline/internal/entity"
)

const (
	ApiVersionV1 = "/v1"

	// Monitoring
	PathHealth    = "/health" // docker health handler
	PathGetRoutes = "/routes"
	PathMetrics   = "/metrics" // prometheus

	// Auth endpoint
	PathAuth              = "/auth"
	PathLogin             = "/login"
	PathToken             = "/token"
	PathCode              = "/codes"
	PathUsersRegistration = "/registration/users"
	PathOrgsRegistration  = "/registration/orgs"

	// User endpoint
	PathUsers      = "/users"
	PathMapOrgs    = "/orgmap"
	PathSearchOrgs = "/search/orgs"

	// Organization endpoint
	PathOrgs       = "/orgs"
	PathServices   = "/services"
	PathTimetables = "/timetables"

	PathWorkers          = "/workers"
	PathWorkersSlots     = PathWorkers + "/slots"
	PathWorkersServices  = PathWorkers + "/services"
	PathWorkersSchedules = PathWorkers + "/schedules"

	// Record endpoint
	PathRecords  = "/records"
	PathFeedback = "/feedbacks"

	// S3 endpoint
	PathMedia = "/media"

	// Analytics service
	PathAnalytics = "/analytics"
)

const (
	authFind      = 5
	usersFind     = 8
	orgsFind      = 15
	recsFind      = 17
	mediaFind     = 18
	analyticsFind = 19
)

var HandlersList = []string{

	PathAuth,
	PathLogin,
	PathUsersRegistration,
	PathOrgsRegistration,
	PathToken,
	PathCode,

	PathUsers,
	PathMapOrgs,
	PathSearchOrgs,

	PathOrgs,
	PathServices,
	PathTimetables,

	PathWorkers,
	PathWorkersSlots,
	PathWorkersServices,
	PathWorkersSchedules,

	PathRecords,
	PathFeedback,

	PathMedia,

	PathAnalytics,
}

type MethodList map[string]string

func newMethodsMap(s *Settings, methods ...string) MethodList {
	m := make(MethodList, len(methods))
	for i := range len(methods) {
		if _, ok := s.SupportedMethodsMap[methods[i]]; ok {
			m[methods[i]] = methods[i]
		}
	}
	return m
}

// Return methods specified by method enum in scope.go
//
// Call endpoint.Methods... - to get all methods
func (mp MethodList) Get(menthods ...string) []string {
	list := make([]string, len(menthods))
	for i := range menthods {
		list[i] = mp[menthods[i]]
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
	// /auth/registration/users
	case PathUsersRegistration:
		mdata.Methods = newMethodsMap(s, http.MethodPost)
		mdata.perms = perms.GrantPermissions(perms.CREATE, perms.NONE)
	// /auth/registration/orgs
	case PathOrgsRegistration:
		mdata.Methods = newMethodsMap(s, http.MethodPost)
		mdata.perms = perms.GrantPermissions(perms.NONE, perms.CREATE)
	// /auth/token
	case PathToken:
		mdata.Methods = newMethodsMap(s, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.CREATE+perms.UPDATE, perms.CREATE+perms.UPDATE)
	// /auth/codes
	case PathCode:
		mdata.Methods = newMethodsMap(s, http.MethodPost, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.CREATE, perms.CREATE)
	// /users
	case PathUsers:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.READ+perms.UPDATE, perms.READ)
	// /users/orgmap
	case PathMapOrgs:
		mdata.Methods = newMethodsMap(s, http.MethodGet)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.NONE)
	// /users/search/orgs
	case PathSearchOrgs:
		mdata.Methods = newMethodsMap(s, http.MethodGet)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.NONE)
	// /orgs  [ Users: Only Read | Orgs: All ]
	case PathOrgs:
		mdata.Methods = newMethodsMap(s, http.MethodGet, http.MethodPut)
		mdata.perms = perms.GrantPermissions(perms.READ, perms.READ+perms.UPDATE)
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
		mdata.Methods = newMethodsMap(s, http.MethodPost, http.MethodDelete, http.MethodGet)
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
	case PathAnalytics:
		mdata.Methods = newMethodsMap(s, http.MethodGet)
		mdata.perms = perms.GrantPermissions(perms.NONE, perms.READ)
	}
	return mdata
}

const (
	ErrPathNotFound = "path not found in %s: \"%s\""
	ErrNoPermission = "%s (%s=%d) has no authorities to call <[%s] %s>"
)

type Routes map[string]endpoint

func NewDefaultRoutes(settings *Settings) Routes {
	r := make(Routes, len(HandlersList))
	for _, path := range HandlersList {
		r[path] = NewEndpointFromPath(settings, path)
	}
	return r
}

// Verifying access rights when entity trying to call a handler
//
// Checks provided method and uri with endpoint's restrictions
func (r Routes) HasAccess(tdata entity.TokenData, uri, method string) error {
	find, prevfind := 0, 0
	switch {
	case strings.Contains(uri, PathAuth):
		find = authFind
	case strings.Contains(uri, PathUsers):
		prevfind = authFind
		find = usersFind
	case strings.Contains(uri, PathOrgs):
		prevfind = usersFind
		find = orgsFind
	case strings.Contains(uri, PathRecords):
		prevfind = orgsFind
		find = recsFind
	case strings.Contains(uri, PathMedia):
		prevfind = recsFind
		find = mediaFind
	case strings.Contains(uri, PathAnalytics):
		prevfind = mediaFind
		find = analyticsFind
	default:
		return fmt.Errorf(ErrPathNotFound, "check prefix", uri)
	}
	ind, isMatched := 0, false
	for i := find; prevfind < i && i <= find; i-- {
		if isMatched = strings.Contains(uri, HandlersList[i]); isMatched {
			ind = i
			break
		}
	}
	if !isMatched {
		return fmt.Errorf(ErrPathNotFound, "handlers list", uri)
	}
	handler, ok := r[HandlersList[ind]]
	if !ok {
		return fmt.Errorf(ErrPathNotFound, "routes", uri)
	}
	if !handler.perms.HasPermission(tdata.IsOrg, method) {
		if tdata.IsOrg {
			return fmt.Errorf(ErrNoPermission, "org", "org_id", tdata.ID, method, uri)
		} else {
			return fmt.Errorf(ErrNoPermission, "user", "user_id", tdata.ID, method, uri)
		}
	}
	return nil
}
