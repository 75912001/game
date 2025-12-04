package role

import (
	"saClient/src/common"
)

type PlayerMgr struct {
	players map[common.UID]*Role
}

func NewPlayerMgr() *PlayerMgr {
	return &PlayerMgr{
		players: make(map[common.UID]*Role),
	}
}
