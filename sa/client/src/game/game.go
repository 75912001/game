package game

import (
	"saClient/src/ui"
)

// GameState 游戏状态
type GameState int

const (
	StateMenu     GameState = iota // 菜单状态
	StateScene                     // 场景中
	StateBattling                  // 战斗中
	StateGameOver                  // 游戏结束
)

type Game struct {
	debugMsg string
	state    GameState // 当前游戏状态

	user *Role
}

func (p *Game) Init() {
	// 初始化为菜单状态
	p.state = StateMenu

	// 添加一个开始按钮
	ui.GUIButtonMgr.AddButton("startBtn", 300, 250, 200, 50, "开始",
		func() (ui.ButtonAction, error) {
			// 点击开始按钮后,切换到游戏状态
			p.state = StateScene
			return ui.ButtonActionDestroy, nil
		},
	)
}
