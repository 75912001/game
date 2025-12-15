package user

import (
	"saClient/src/res"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	bottomCenterWX float32 // 脚底中心 World 坐标 X
	bottomCenterWY float32 // 脚底中心 World 坐标 Y

	bottomCenterTX float32 // 脚底中心 Tile 坐标 X
	bottomCenterTY float32 // 脚底中心 Tile 坐标 Y

	// 角色中心 World 坐标 (用于渲染和摄像机跟随)
	centerWX float32
	centerWY float32

	orientation     uint32               // 方向
	image           *ebitenv2.Image      // 角色图片
	roleImageSprite *res.RoleImageSprite // 角色图片配置
}
