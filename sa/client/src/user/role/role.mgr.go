package role

import (
	xmap "github.com/75912001/xlib/map"
	"saClient/src/common"
)

type RoleMgr struct {
	Roles *xmap.MapMgr[common.RoleUUID, *Role]
}

func NewRoleMgr() *RoleMgr {
	return &RoleMgr{
		Roles: xmap.NewMapMgr[common.RoleUUID, *Role](),
	}
}
