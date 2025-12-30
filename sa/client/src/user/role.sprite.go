package user

import (
	"saClient/src/proto"
	"saClient/src/res"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// RoleSprite 角色-精灵
type RoleSprite struct {
	wx float32 // 脚底中心 World 坐标 X
	wy float32 // 脚底中心 World 坐标 Y

	tx float32 // 脚底中心 Tile 坐标 X
	ty float32 // 脚底中心 Tile 坐标 Y

	centerWX float32 // 角色中心 World 坐标 X
	centerWY float32 // 角色中心 World 坐标 Y

	action          proto.RoleAction       // 当前播放的动作 (Move, AttackAxe等)
	orientation     proto.AssetOrientation // 方向
	image           *ebitenv2.Image        // 角色图片
	roleImageSprite *res.RoleImageSprite   // 角色图片配置
}
