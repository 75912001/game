package user

import (
	"saClient/src/common"
	"saClient/src/user/role"
)

func (p *User) Update() {
	p.roleMgr.Roles.Foreach(
		func(key common.RoleUUID, role *role.Role) (isContinue bool) {
			role.Update()
			return true
		},
	)
}
