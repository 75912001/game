package user

import (
	"saClient/src/common"
)

func (p *User) Update() {
	p.roleMgr.Roles.Foreach(
		func(key common.RoleUUID, role *Role) (isContinue bool) {
			role.Update()
			return true
		},
	)
}
