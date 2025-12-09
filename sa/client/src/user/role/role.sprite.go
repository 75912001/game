package role

import "saClient/src/res"

// RoleSprite 角色-精灵
type RoleSprite struct {
	x               int
	y               int
	direction       int
	roleImageSprite *res.RoleImageSprite
}
