package game

import (
	"saClient/src/role"
	"saClient/src/ui"
)

type Game struct {
	state State // 当前游戏状态
	role  *role.Role
}

func (p *Game) Init() {
	// 初始化为菜单状态
	p.state = State_StartMenu

	// 添加一个开始按钮
	ui.GUIButtonMgr.AddButton("startBtn", 300, 250, 200, 50, "开始",
		func() (ui.ButtonAction, error) {
			// 点击开始按钮后,切换到游戏状态
			p.state = State_Scene
			// 销毁按钮
			return ui.ButtonActionDestroy, nil
		},
	)
}
