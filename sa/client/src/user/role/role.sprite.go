package role

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/res"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	x               int
	y               int
	direction       int
	image           *ebitenv2.Image
	roleImageSprite *res.RoleImageSprite
}
