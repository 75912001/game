package ui

import (
	ebitenui "github.com/ebitenui/ebitenui"
	ebitenuiimage "github.com/ebitenui/ebitenui/image"
	ebitenuiwidget "github.com/ebitenui/ebitenui/widget"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	imagecolor "image/color"
)

var GUIMgr *UIMgr

func init() {
	GUIMgr = NewUIMgr()
}

type UIMgr struct {
	ui       *ebitenui.UI // UI实例
	debugMsg string       // 调试信息
}

func NewUIMgr() *UIMgr {
	// 创建根容器,不使用布局管理器(允许手动定位)
	return &UIMgr{
		ui: &ebitenui.UI{
			Container: ebitenuiwidget.NewContainer(
				ebitenuiwidget.ContainerOpts.BackgroundImage(ebitenuiimage.NewNineSliceColor(imagecolor.Transparent)),
			),
		},
	}
}

// AddButton 添加一个按钮
func (p *UIMgr) AddButton(button *ebitenuiwidget.Button) {
	p.ui.Container.AddChild(button)
}

// RemoveButton 移除一个按钮
func (p *UIMgr) RemoveButton(button *ebitenuiwidget.Button) {
	p.ui.Container.RemoveChild(button)
}

// Update 更新 UI
func (p *UIMgr) Update() {
	p.ui.Update()
}

// Draw 绘制 UI
func (p *UIMgr) Draw(screen *ebitenv2.Image) {
	p.ui.Draw(screen)
	if p.debugMsg != "" {
		Printf(screen, 0, 0, p.debugMsg)
	}
}

func (p *UIMgr) SetDebugMsg(msg string) {
	p.debugMsg = msg
}
