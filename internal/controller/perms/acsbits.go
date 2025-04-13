package perms

import (
	"net/http"
)

const (
	CREATE = 1 << iota
	READ
	UPDATE
	DELETE
	ALL  = CREATE + READ + UPDATE + DELETE
	NONE = 0
)

const (
	halfbyte = 4
	ORGMASK  = 15
	USERMASK = ORGMASK << halfbyte
)

type PermissionBits uint8

// designed as UNIX permission bits
//
// user-org
//
// xxxx-xxxx - permission bits
//
//		Examples:
//		1. CRUD-CRUD = 1515 - ALL-ALL
//		2. XRXX-CRUD = 0215 - READ-ALL
//	    ...
func GrantPermissions(user, org uint8) PermissionBits {
	bits := user
	bits = bits << 4
	bits |= org
	return PermissionBits(bits)
}

func (pb PermissionBits) HasPermission(isOrg bool, method string) bool {
	var perms PermissionBits
	if isOrg {
		perms = pb & ORGMASK
	} else {
		perms = pb & USERMASK >> halfbyte
	}
	switch method {
	case http.MethodGet:
		return perms&READ != 0
	case http.MethodPut:
		return perms&UPDATE != 0
	case http.MethodPost:
		return perms&CREATE != 0
	case http.MethodDelete:
		return perms&DELETE != 0
	case http.MethodPatch:
		return perms&UPDATE != 0
	default:
		return false
	}
}
