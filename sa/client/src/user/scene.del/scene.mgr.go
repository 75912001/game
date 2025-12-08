package scene_del

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/common"
)

// todo menglc 场景中包含多个 sheet, 如 excel 一样,底部有多个标签页可以选择. 每个 sheet 包含不同的角色游戏
//  通过用户使用按键 1,2,3,4,5...切换不同的角色

type SceneMgr struct {
	scenes []*Scene //
}

func NewSceneMgr() *SceneMgr {
	sceneMgr := &SceneMgr{
		scenes: make([]*Scene, 0, 3),
	}
}

func (p *SceneMgr) Update() error {
	return nil
}

func (p *SceneMgr) Draw(screen *ebitenv2.Image) {

}
