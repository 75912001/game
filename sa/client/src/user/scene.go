package user

import (
	"saClient/src/arpg"
	"saClient/src/common"
)

type Scene struct {
	_map         *Map               // Tiled 地图
	transition   *SceneTransition   // 场景过渡效果
	spawnManager *arpg.SpawnManager // 刷怪点管理器
}

func NewScene(mapID common.AssetID) *Scene {
	scene := &Scene{
		_map: NewMap(mapID),
	}
	scene.transition = newSceneTransition(scene._map.mapCfg.SceneTransition)
	scene.spawnManager = arpg.NewSpawnManager(scene._map.tiledMapCfg)
	return scene
}

func (p *Scene) Update() {
	p._map.Update()
	p.spawnManager.Update()
}

// GetMapTileSize 获取地图 tile 尺寸
func (p *Scene) GetMapTileSize() (width, height int) {
	return p._map.tiledMapCfg.Width, p._map.tiledMapCfg.Height
}
