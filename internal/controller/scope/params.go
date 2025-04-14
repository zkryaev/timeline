package scope

// supported query param
const (
	// Entity_id
	USER_ID    = "user_id"
	ORG_ID     = "org_id"
	WORKER_ID  = "worker_id"
	SERVICE_ID = "service_id"
	RECORD_ID  = "record_id"

	// Pagination
	LIMIT = "limit"
	PAGE  = "page"

	// Limitations
	WEEKDAY = "weekday"
	AS_LIST = "as_list"
	FRESH   = "fresh"

	// SEARCH
	NAME    = "name"
	TYPE    = "type"
	SORT_BY = "sort_by"
	ORDER   = "order"

	// MAP
	MIN_LAT = "min_lat"
	MIN_LON = "min_lon"
	MAX_LAT = "max_lat"
	MAX_LON = "max_lon"
)

const (
	SINGLE = false
	LIST   = true
)

type SupportedParams map[string]map[string]struct{}

func defaultSupportedParams() SupportedParams {
	return SupportedParams{
		INT:     {USER_ID: {}, ORG_ID: {}, WORKER_ID: {}, SERVICE_ID: {}, RECORD_ID: {}, LIMIT: {}, PAGE: {}, WEEKDAY: {}},
		BOOL:    {AS_LIST: {}, FRESH: {}},
		STRING:  {NAME: {}, TYPE: {}, SORT_BY: {}, ORDER: {}},
		FLOAT32: {MIN_LAT: {}, MIN_LON: {}, MAX_LAT: {}, MAX_LON: {}},
	}
}

// if name found: param_type
// else: ""
func (sp SupportedParams) GetParam(name string) string {
	for ptype, pmap := range sp {
		_, ok := pmap[name]
		if ok {
			return ptype
		}
	}
	return ""
}
