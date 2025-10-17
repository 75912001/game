package ui

import (
	"github.com/ebitenui/ebitenui"
	ebitenuiimage "github.com/ebitenui/ebitenui/image"
	ebitenuiwidget "github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	imgcolor "image/color"
)

var GUIMgr *UIMgr

func init() {
	GUIMgr = NewUIMgr()
}

type UIMgr struct {
	ui *ebitenui.UI // UI实例
}

func NewUIMgr() *UIMgr {
	// 创建根容器，不使用布局管理器（允许手动定位）
	return &UIMgr{
		ui: &ebitenui.UI{
			Container: ebitenuiwidget.NewContainer(
				ebitenuiwidget.ContainerOpts.BackgroundImage(ebitenuiimage.NewNineSliceColor(imgcolor.Transparent)),
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
func (p *UIMgr) Draw(screen *ebiten.Image) {
	p.ui.Draw(screen)
}
