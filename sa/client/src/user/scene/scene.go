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
	sceneMap      *Map           // 地图
	plantMgr      *PlantMgr      // 植物-管理器
	itemMgr       *ItemMgr       // 物品-管理器
}

func NewScene(mapID common.AssetID) *Scene {
	scene := &Scene{
		buildingMgr:   NewBuildingMgr(),
		decorationMgr: NewDecorationMgr(),
		sceneMap:      NewMap(mapID),
		plantMgr:      NewPlantMgr(),
		itemMgr:       NewItemMgr(),
	}
	return scene
}

func (p *Scene) Update() {
	p.sceneMap.Update()
	p.buildingMgr.Update()
	p.plantMgr.Update()
	p.decorationMgr.Update()
	p.itemMgr.Update()
}

// GetMapSize 获取地图尺寸
func (p *Scene) GetMapSize() (width, height int) {
	return p.sceneMap.cfg.Width, p.sceneMap.cfg.Height
}

func (p *Scene) Draw(screen *ebitenv2.Image, camera *camera.Camera) {
	// 填充草地一样的绿色背景
	screen.Fill(color.RGBA{
		R: 34,
		G: 139,
		B: 34,
		A: 255,
	})

	p.sceneMap.Draw(screen, camera)
	p.buildingMgr.Draw(screen)
	p.plantMgr.Draw(screen)
	p.decorationMgr.Draw(screen)
	p.itemMgr.Draw(screen)

	// 根据游戏状态绘制不同的内容
	// switch p.state {
	// case scene.State_StartMenu: // 开始菜单状态下绘制 UI
	//	ui.GUIMgr.Draw(screen)
	// case scene.State_Scene: // 场景状态,绘制场景背景等
	// case scene.State_Battling: // 战斗状态,绘制战斗相关内容
	// case scene.State_GameOver: // 游戏结束,显示游戏结束信息
	//	ui.Printf(screen, 280, 280, "*** 游戏结束 ***")
	// }
}
