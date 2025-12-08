package game

import (
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
	// 添加一个登录按钮
	ui.GUIButtonMgr.AddButton("loginBtn", 300, 250, 200, 50, "登录",
		func() (ui.ButtonAction, error) { // 点击开始按钮后,切换到场景
			err := p.user.Login("", "")
			if err != nil {
				return ui.ButtonActionShow, nil
			}
			return ui.ButtonActionDestroy, nil // 销毁按钮
		},
	)
}
