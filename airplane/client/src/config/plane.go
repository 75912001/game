package config

import (
	"airplaneClient/src/common"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

type Plane struct {
	common.PlaneKey
	Width       uint32                                     // 宽度
	Height      uint32                                     // 高度
	PartsFrames [common.PlanePartTypeMax][]*ebitenv2.Image // 各部件-动画帧
}
