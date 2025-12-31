package res

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

// PetAttack 宠物攻击动画数据
type PetAttack struct {
	Frames    [proto.AssetOrientation_AssetOrientation_Max][]*ebitenv2.Image // 动画帧
	FrameInfo [proto.AssetOrientation_AssetOrientation_Max][]*PetImageSprite // 动画帧信息
}

// NewPetAttack 创建宠物攻击动画数据
func NewPetAttack() *PetAttack {
	return &PetAttack{}
}
