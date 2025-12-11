package scene

import (
	"image/color"
	"saClient/src/common"
	"saClient/src/proto"
	"saClient/src/res/tiled"
	"saClient/src/user/camera"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

type Scene struct {
	buildingMgr   *BuildingMgr   // 建筑-管理器
	decorationMgr *DecorationMgr // 装饰物-管理器
	sceneMap      *Map           // 普通地图
	tiledMap      *TiledMap      // Tiled 地图
	plantMgr      *PlantMgr      // 植物-管理器
	itemMgr       *ItemMgr       // 物品-管理器
}

func NewScene(mapID common.AssetID) *Scene {
	scene := &Scene{
		buildingMgr:   NewBuildingMgr(),
		decorationMgr: NewDecorationMgr(),
		plantMgr:      NewPlantMgr(),
		itemMgr:       NewItemMgr(),
	}

	// 优先检查是否有 Tiled 地图
	if tiled.GMapMgr.Maps.Get(mapID) != nil {
		scene.tiledMap = NewTiledMap(mapID)
	} else {
		scene.sceneMap = NewMap(mapID)
	}
	return scene
}

func (p *Scene) Update() {
	switch p.GetResMapType() {
	case proto.ResMapType_ResMapType_Normal:
		p.sceneMap.Update()
	case proto.ResMapType_ResMapType_Tiled:
		p.tiledMap.Update()
	}
	p.buildingMgr.Update()
	p.plantMgr.Update()
	p.decorationMgr.Update()
	p.itemMgr.Update()
}

// GetResMapType 获取-资源-地图-类型
func (p *Scene) GetResMapType() proto.ResMapType {
	if p.sceneMap != nil {
		return proto.ResMapType_ResMapType_Normal
	}
	if p.tiledMap != nil {
		return proto.ResMapType_ResMapType_Tiled
	}
	return proto.ResMapType_ResMapType_Unknow
}

// GetMapSize 获取地图尺寸
func (p *Scene) GetMapSize() (width, height int) {
	switch p.GetResMapType() {
	case proto.ResMapType_ResMapType_Normal:
		return p.sceneMap.cfg.Width, p.sceneMap.cfg.Height
	case proto.ResMapType_ResMapType_Tiled:
		return p.tiledMap.GetMapSize()
	}
	return 0, 0
}

func (p *Scene) Draw(screen *ebitenv2.Image, camera *camera.Camera) {
	// 填充草地一样的绿色背景
	screen.Fill(color.RGBA{
		R: 34,
		G: 139,
		B: 34,
		A: 255,
	})

	// 根据地图类型绘制
	switch p.GetResMapType() {
	case proto.ResMapType_ResMapType_Normal:
		p.sceneMap.Draw(screen, camera)
	case proto.ResMapType_ResMapType_Tiled:
		p.tiledMap.Draw(screen, camera)
	}

	p.buildingMgr.Draw(screen)
	p.plantMgr.Draw(screen)
	p.decorationMgr.Draw(screen)
	p.itemMgr.Draw(screen)
}
