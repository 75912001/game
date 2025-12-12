package role

import (
	"saClient/src/res"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	centerTX int // 角色中心坐标-tx
	centerTY int // 角色中心Y坐标-ty

	bottomCenterTX int // 底部中心点坐标-tx
	bottomCenterTY int // 底部中心点坐标-ty

	orientation     int                  // 方向
	image           *ebitenv2.Image      // 角色图片
	roleImageSprite *res.RoleImageSprite // 角色图片配置
}
