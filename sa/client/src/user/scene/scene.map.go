package scene

import (
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/user/camera"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// Map 场景-地图
type Map struct {
	id  common.AssetID // 地图ID
	cfg *cfg.Map       // 地图配置
}

func NewMap(mapID common.AssetID) *Map {
	m := &Map{
		id: mapID,
	}
	// 从配置管理器获取地图配置
	m.cfg = cfg.GMapMgr.Maps.Get(mapID)
	return m
}

func (p *Map) Update() {

}

func (p *Map) Draw(screen *ebitenv2.Image, camera *camera.Camera) {
	op := &ebitenv2.DrawImageOptions{}
	// 将地图图片绘制到屏幕上
	op.GeoM.Translate(float64(-camera.ScreenX), float64(-camera.ScreenY))
	screen.DrawImage(p.cfg.ResMap.Image, op)
}
