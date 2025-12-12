package role

import (
	"saClient/src/res"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	// 脚底中心 World 坐标 (float64 用于平滑移动)
	bottomCenterWorldX float64
	bottomCenterWorldY float64

	// 角色中心 World 坐标 (用于渲染和摄像机跟随)
	centerWorldX float64
	centerWorldY float64

	orientation     int                  // 方向
	image           *ebitenv2.Image      // 角色图片
	roleImageSprite *res.RoleImageSprite // 角色图片配置
}
