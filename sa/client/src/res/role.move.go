package res

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

type RoleMove struct {
	Frames    [proto.AssetDirection_AssetDirection_Max][]*ebitenv2.Image  // 动画帧
	FrameInfo [proto.AssetDirection_AssetDirection_Max][]*RoleImageSprite // 向上-动画帧信息 配置表中 (镜像帧信息与原帧相同)
}

func NewRoleMove() *RoleMove {
	return &RoleMove{}
}
