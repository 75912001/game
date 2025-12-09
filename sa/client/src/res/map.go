package res

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/common"
)

// Map 地图资源
type Map struct {
	ID    common.AssetID  // 地图ID
	Image *ebitenv2.Image // 地图图片
}
