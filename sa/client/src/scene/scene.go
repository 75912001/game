package scene

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/role"
)

var GScene *Scene // 全局-当前场景

type Scene struct {
	role      *role.Role      // 角色
	playerMgr *role.PlayerMgr // 其他玩家-管理器

	buildingMgr   *BuildingMgr   // 建筑-管理器
	decorationMgr *DecorationMgr // 装饰物-管理器
	sceneMap      *Map           // 地图
	plantMgr      *PlantMgr      // 植物-管理器
	itemMgr       *ItemMgr       // 物品-管理器
}

type Mgr struct {
	scenes map[int32]*Scene
}

func (p *Scene) Update() error {
	return nil
}

func (p *Scene) Draw(screen *ebitenv2.Image) {
}
