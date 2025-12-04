package role

import (
	"saClient/src/common"
)

type Mgr struct {
	players map[common.UID]*Role
}

func NewMgr() *Mgr {
	return &Mgr{
		players: make(map[common.UID]*Role),
	}
}
