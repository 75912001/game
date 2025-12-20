package res

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

type PetMove struct {
	Frames    [proto.AssetOrientation_AssetOrientation_Max][]*ebitenv2.Image // 动画帧
	FrameInfo [proto.AssetOrientation_AssetOrientation_Max][]*PetImageSprite // 动画帧信息 配置表中 (镜像帧信息与原帧相同)
}

func NewPetMove() *PetMove {
	return &PetMove{}
}
