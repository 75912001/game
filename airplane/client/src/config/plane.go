package config

import (
	"airplaneClient/src/common"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

type Plane struct {
	common.PlaneKey
	PartsFramesData [common.PlanePartTypeMax][]*ebitenv2.Image  // 各部件-动画帧
	PartsFramesInfo [common.PlanePartTypeMax][]*ConfigPlanePart // 各部件-动画帧信息
}
