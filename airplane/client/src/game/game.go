package game

import (
	"airplaneClient/src/ui"
)

type Game struct {
	debugMsg   string
	clickCount int
}

func (g *Game) Init() {
	// 添加一个开始按钮
	ui.GUIButtonMgr.AddButton("startBtn", 300, 250, 200, 50, "开始", func() (ui.ButtonAction, error) {
		return ui.ButtonActionDestroy, nil
	})
}
