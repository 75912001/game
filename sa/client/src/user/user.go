package user

import (
	"saClient/src/common"
	"saClient/src/proto"
)

type User struct {
	record *proto.UserRecord
}

func NewUser() *User {
	return &User{}
}

func (p *User) Login(account string, password string) error {
	p.record = &proto.UserRecord{}
	// 初始化用户数据 todo menglc 模拟数据
	return nil
}

func (p *User) SelectRole(roleUUID common.RoleUUID) error {
	return nil
}
