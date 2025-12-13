package user

import (
	"image/color"
	"saClient/src/common"
	"saClient/src/user/camera"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

type Scene struct {
	buildingMgr   *BuildingMgr   // 建筑-管理器
	decorationMgr *DecorationMgr // 装饰物-管理器
	_map          *Map           // Tiled 地图
	plantMgr      *PlantMgr      // 植物-管理器
	itemMgr       *ItemMgr       // 物品-管理器
}

func NewScene(mapID common.AssetID) *Scene {
	scene := &Scene{
		buildingMgr:   NewBuildingMgr(),
		decorationMgr: NewDecorationMgr(),
		_map:          NewMap(mapID),
		plantMgr:      NewPlantMgr(),
		itemMgr:       NewItemMgr(),
	}
	return scene
}

func (p *Scene) Update() {
	p._map.Update()
	p.buildingMgr.Update()
	p.plantMgr.Update()
	p.decorationMgr.Update()
	p.itemMgr.Update()
}

// GetMapPixeSize 获取地图像素尺寸
func (p *Scene) GetMapPixeSize() (width, height int) {
	return p._map.cfg.PixelW, p._map.cfg.PixelH
}

// GetMapID 获取地图ID
func (p *Scene) GetMapID() common.AssetID {
	return p._map.cfg.ID
}

// GetMapTileSize 获取地图 tile 尺寸
func (p *Scene) GetMapTileSize() (width, height int) {
	return p._map.cfg.Width, p._map.cfg.Height
}

// ClampTileBounds 将 tile 坐标限制在地图边界内
func (p *Scene) ClampTileBounds(tileX, tileY float32) (clampedTX, clampedTY float32) {
	return p._map.ClampTileBounds(tileX, tileY)
}

// TileToWorld tile 坐标转换为世界坐标
func (p *Scene) TileToWorld(tileX, tileY float32) (worldX, worldY float32) {
	return p._map.cfg.IsometricCT.T2W(tileX, tileY)
}

// WorldToTile 世界坐标转换为 tile 坐标
func (p *Scene) WorldToTile(worldX, worldY float32) (tileX, tileY float32) {
	return p._map.cfg.IsometricCT.W2T(worldX, worldY)
}

func (p *Scene) Draw(screen *ebitenv2.Image, camera *camera.Camera) {
	// 填充草地一样的绿色背景
	screen.Fill(color.RGBA{
		R: 34,
		G: 139,
		B: 34,
		A: 255,
	})

	p._map.Draw(screen, camera)
	p.buildingMgr.Draw(screen)
	p.plantMgr.Draw(screen)
	p.decorationMgr.Draw(screen)
	p.itemMgr.Draw(screen)
}
