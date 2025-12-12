package scene

import (
	"image/color"
	"saClient/src/common"
	"saClient/src/user/camera"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

type Scene struct {
	buildingMgr   *BuildingMgr   // 建筑-管理器
	decorationMgr *DecorationMgr // 装饰物-管理器
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

	scene.tiledMap = NewTiledMap(mapID)
	return scene
}

func (p *Scene) Update() {
	p.tiledMap.Update()
	p.buildingMgr.Update()
	p.plantMgr.Update()
	p.decorationMgr.Update()
	p.itemMgr.Update()
}

// GetMapSize 获取地图尺寸
func (p *Scene) GetMapSize() (width, height int) {
	return p.tiledMap.GetMapSize()
}

// ClampToMapBounds 将坐标限制在地图边界内
// 对于 Tiled 等距地图，边界是菱形；对于普通地图，边界是矩形
func (p *Scene) ClampToMapBounds(worldX, worldY float64) (clampedX, clampedY float64) {
	return p.tiledMap.ClampToMapBounds(worldX, worldY)
}

// IsInMapBounds 检查坐标是否在地图边界内
func (p *Scene) IsInMapBounds(worldX, worldY float64) bool {
	return p.tiledMap.IsInMapBounds(worldX, worldY)
}

func (p *Scene) Draw(screen *ebitenv2.Image, camera *camera.Camera) {
	// 填充草地一样的绿色背景
	screen.Fill(color.RGBA{
		R: 34,
		G: 139,
		B: 34,
		A: 255,
	})

	p.tiledMap.Draw(screen, camera)
	p.buildingMgr.Draw(screen)
	p.plantMgr.Draw(screen)
	p.decorationMgr.Draw(screen)
	p.itemMgr.Draw(screen)
}
