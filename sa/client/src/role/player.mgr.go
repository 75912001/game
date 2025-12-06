package role

import (
	xmap "github.com/75912001/xlib/map"
	"saClient/src/common"
)

type PlayerMgr struct {
	players *xmap.MapMgr[common.RoleID, *Role]
}

func NewPlayerMgr() *PlayerMgr {
	return &PlayerMgr{
		players: xmap.NewMapMgr[common.RoleID, *Role](),
	}
}
