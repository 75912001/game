package user

import (
	"saClient/src/common"
)

type Scene struct {
	buildingMgr   *BuildingMgr     // 建筑-管理器
	decorationMgr *DecorationMgr   // 装饰物-管理器
	_map          *Map             // Tiled 地图
	plantMgr      *PlantMgr        // 植物-管理器
	itemMgr       *ItemMgr         // 物品-管理器
	transition    *SceneTransition // 场景过渡效果
}

func NewScene(mapID common.AssetID) *Scene {
	scene := &Scene{
		buildingMgr:   NewBuildingMgr(),
		decorationMgr: NewDecorationMgr(),
		_map:          NewMap(mapID),
		plantMgr:      NewPlantMgr(),
		itemMgr:       NewItemMgr(),
	}
	scene.transition = newSceneTransition(scene._map.mapCfg.SceneTransition)
	return scene
}

func (p *Scene) Update() {
	p._map.Update()
	p.buildingMgr.Update()
	p.plantMgr.Update()
	p.decorationMgr.Update()
	p.itemMgr.Update()
}

// GetMapTileSize 获取地图 tile 尺寸
func (p *Scene) GetMapTileSize() (width, height int) {
	return p._map.tiledMapCfg.Width, p._map.tiledMapCfg.Height
}
