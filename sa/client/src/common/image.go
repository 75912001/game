package common

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

type Image struct {
	resType proto.ResType

	resMajorID uint32 // 资源大类ID
	resMinorID uint32 // 资源小类ID

	width  uint32 // 宽度
	height uint32 // 高度

	frames []*ebitenv2.Image // 动画帧
}
