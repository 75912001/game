package common

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

type Image struct {
	assetType proto.AssetType // 资源类型
	assetID   AssetID         // 资源ID
	*proto.Image

	frames []*ebitenv2.Image // 动画帧
}
