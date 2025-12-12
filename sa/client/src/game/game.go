package game

import (
	"saClient/src/cfg"
	"saClient/src/ui"
	"saClient/src/user"
)

type Game struct {
	user *user.User // 当前用户
}

func NewGame() *Game {
	return &Game{
		user: user.NewUser(),
	}
}

func (p *Game) Init() {
	// 按钮尺寸
	btnWidth := 200
	btnHeight := 50

	// 计算居中位置
	btnX := (cfg.GCommon.ScreenMaxWidth - btnWidth) / 2
	btnY := (cfg.GCommon.ScreenMaxHeight - btnHeight) / 2

	// 添加一个登录按钮（居中显示）
	ui.GUIButtonMgr.AddButton("loginBtn", btnX, btnY, btnWidth, btnHeight, "登录",
		func() (ui.ButtonAction, error) { // 点击开始按钮后,切换到场景
			err := p.user.Login("", "")
			if err != nil {
				return ui.ButtonActionShow, nil
			}
			return ui.ButtonActionDestroy, nil // 销毁按钮
		},
	)
}
