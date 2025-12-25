package user

import (
	"saClient/src/proto"
	"saClient/src/res"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	actionAnchorPointWX float32 // 动作锚点-脚底中心 World 坐标 X
	actionAnchorPointWY float32 // 动作锚点-脚底中心 World 坐标 Y

	actionAnchorPointTX float32 // 动作锚点-脚底中心 Tile 坐标 X
	actionAnchorPointTY float32 // 动作锚点-脚底中心 Tile 坐标 Y

	// 角色中心 World 坐标 (用于渲染和摄像机跟随)
	cameraAnchorPointWX float32 // camera锚点-角色中心 World 坐标 X
	cameraAnchorPointWY float32 // camera锚点-角色中心 World 坐标 Y

	orientation     proto.AssetOrientation // 方向
	image           *ebitenv2.Image        // 角色图片
	roleImageSprite *res.RoleImageSprite   // 角色图片配置
}
