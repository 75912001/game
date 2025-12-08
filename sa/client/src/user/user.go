package user

import (
	"saClient/src/common"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (p *User) SelectRole(roleUUID common.RoleUUID) error {
	return nil
}
