package user

import (
	"image/color"
	"saClient/src/cfg"
	"saClient/src/common"
	"saClient/src/user/camera"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
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
	scene.transition = newSceneTransition(scene._map.cfgMap.SceneTransition)
	return scene
}

func (p *Scene) Update() {
	p._map.Update()
	p.buildingMgr.Update()
	p.plantMgr.Update()
	p.decorationMgr.Update()
	p.itemMgr.Update()
}

// TransitionIn 转场进入 (黑幕展开，显示场景)
func (p *Scene) TransitionIn() {
	p.transition.TransitionIn()
}

// TransitionOut 转场退出 (使用指定配置)
func (p *Scene) TransitionOut(transCfg *cfg.SceneTransition) {
	p.transition.SetConfig(transCfg)
	p.transition.TransitionOut()
}

// UpdateTransition 更新过渡状态
// 返回值:
//   - isActive: 是否正在过渡中
//   - isCloseComplete: 转场退出完成
//   - isOpenComplete: 转场进入完成
func (p *Scene) UpdateTransition() (isActive bool, isCloseComplete bool, isOpenComplete bool) {
	return p.transition.Update()
}

// IsTransitionActive 是否正在过渡中
func (p *Scene) IsTransitionActive() bool {
	return p.transition.IsActive()
}

// GetMapTileSize 获取地图 tile 尺寸
func (p *Scene) GetMapTileSize() (width, height int) {
	return p._map.cfg.Width, p._map.cfg.Height
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

	// 绘制过渡效果 (最上层)
	p.transition.Draw(screen)
}
