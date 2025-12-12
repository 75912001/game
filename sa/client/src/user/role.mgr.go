package user

import (
	"saClient/src/common"

	xmap "github.com/75912001/xlib/map"
)

type RoleMgr struct {
	Roles *xmap.MapMgr[common.RoleUUID, *Role]
}

func NewRoleMgr() *RoleMgr {
	return &RoleMgr{
		Roles: xmap.NewMapMgr[common.RoleUUID, *Role](),
	}
}
