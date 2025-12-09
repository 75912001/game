package scene

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/cfg"
	"saClient/src/common"
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

func (p *Map) Draw(screen *ebitenv2.Image, cameraX, cameraY int) {
	op := &ebitenv2.DrawImageOptions{}
	// 相机能拍到的区域是以屏幕中心为基准的
	// 计算左上角的位置
	screenX := cameraX - common.ScreenWidth/2
	screenY := cameraY - common.ScreenHeight/2
	// 将地图图片绘制到屏幕上
	op.GeoM.Translate(float64(-screenX), float64(-screenY))
	screen.DrawImage(p.cfg.ResMap.Image, op)
}
