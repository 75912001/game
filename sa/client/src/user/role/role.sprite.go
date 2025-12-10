package role

import (
	"saClient/src/res"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	centerX int // 角色中心坐标-x
	centerY int // 角色中心Y坐标-y

	bottomCenterX int // 底部中心点坐标-x
	bottomCenterY int // 底部中心点坐标-y

	orientation     int                  // 方向
	image           *ebitenv2.Image      // 角色图片
	roleImageSprite *res.RoleImageSprite // 角色图片配置
}
