package role

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/res"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	x               int
	y               int
	direction       int                  // 方向
	image           *ebitenv2.Image      // 角色图片
	roleImageSprite *res.RoleImageSprite // 角色图片配置
}
