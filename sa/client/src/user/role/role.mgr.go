package role

import (
	xmap "github.com/75912001/xlib/map"
	"saClient/src/common"
)

type RoleMgr struct {
	m *xmap.MapMgr[common.RoleUUID, *Role]
}

func NewRoleMgr() *RoleMgr {
	return &RoleMgr{
		m: xmap.NewMapMgr[common.RoleUUID, *Role](),
	}
}
