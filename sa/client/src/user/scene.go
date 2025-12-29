package user

import (
	"saClient/src/common"
)

type Scene struct {
	_map       *Map             // Tiled 地图
	transition *SceneTransition // 场景过渡效果
}

func NewScene(mapID common.AssetID) *Scene {
	scene := &Scene{
		_map: NewMap(mapID),
	}
	scene.transition = newSceneTransition(scene._map.mapCfg.SceneTransition)
	return scene
}

func (p *Scene) Update() {
	p._map.Update()
}

// GetEnemies 获取场景中的所有敌人
func (p *Scene) GetEnemies() []*ArpgEnemy {
	return p._map.spawnManager.GetAllEnemies()
}

// GetMapTileSize 获取地图 tile 尺寸
func (p *Scene) GetMapTileSize() (width, height int) {
	return p._map.tiledMapCfg.TileCountW, p._map.tiledMapCfg.TileCountH
}
