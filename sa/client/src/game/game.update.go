package game

import (
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
	"saClient/src/ui"
)

func (p *Game) Update() error {
	ui.GUIMgr.Update()

	// 根据游戏状态执行不同的更新逻辑
	switch p.state {
	case State_StartMenu: // 菜单状态下的更新逻辑
		p.updateMenu()
	case State_Scene: // 场景中状态下的更新逻辑
		p.updateScene()
	case State_Battling: // 游戏进行中的更新逻辑
		p.updateBattling()
	case State_GameOver: // 游戏结束状态下的更新逻辑
		p.updateGameOver()
	}
	return nil
}

// updateMenu 菜单状态的更新
func (p *Game) updateMenu() {
	ui.GUIMgr.SetDebugMsg("菜单状态 - 点击开始按钮开始游戏")
}

func (p *Game) updateScene() {
}

// updateBattling 游戏进行中的更新
func (p *Game) updateBattling() {
	if p.scene == nil {
		p.startGame()
	}

	if p.scene != nil {
		_ = p.scene.Update()
		if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeyW) { // 用户按w
		}
		if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeyS) { // 用户按s
		}
		if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeyA) { // 用户按a
		}
		if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeyD) { // 用户按d
		}
		if ebitenv2inpututil.IsKeyJustPressed(ebitenv2.KeySpace) { // 用户按空格键
		}
	}
	ui.GUIMgr.SetDebugMsg(fmt.Sprintf("游戏进行中"))
}

// startGame 开始游戏,初始化游戏对象
func (p *Game) startGame() {
	// 创建玩家角色
	// p.scene = scene.role.NewRole(nil)
}

// updateGameOver 游戏结束状态的更新
func (p *Game) updateGameOver() {
	ui.GUIMgr.SetDebugMsg("游戏结束")
}
