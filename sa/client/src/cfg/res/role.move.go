package res

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

type RoleMove struct {
	Frames    [proto.RoleDirection_RoleDirection_Max]*ebitenv2.Image    // 动画帧
	FrameInfo [proto.RoleDirection_RoleDirection_Max][]*RoleImageSprite // 向上-动画帧信息 配置表中
}

func NewRoleMove() *RoleMove {
	return &RoleMove{}
}
