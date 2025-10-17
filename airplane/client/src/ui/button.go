package ui

import (
	apcfont "airplaneClient/src/resources/font"
	ebitenuiimage "github.com/ebitenui/ebitenui/image"
	ebitenuiwidget "github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	imgcolor "image/color"
	"log"
)

var GUIButtonMgr *ButtonMgr

// 按钮动作
const (
	ButtonActionNone    ButtonAction = iota // 无操作
	ButtonActionHide                        // 隐藏
	ButtonActionShow                        // 显示
	ButtonActionDestroy                     // 销毁
)

type ButtonAction uint32

func init() {
	GUIButtonMgr = NewButtonMgr()
}

// ButtonMgr 管理所有按钮
type ButtonMgr struct {
	buttons map[string]*ButtonWrapper // 按钮映射 key: 按钮ID, value: 按钮包装
}

// ButtonWrapper 包装 ebitenuiwidget 的按钮
type ButtonWrapper struct {
	id       string                       // 按钮ID
	button   *ebitenuiwidget.Button       // 按钮
	visible  bool                         // 是否可见
	callback func() (ButtonAction, error) // 回调函数
}

// NewButtonMgr 创建按钮管理器
func NewButtonMgr() *ButtonMgr {
	return &ButtonMgr{
		buttons: make(map[string]*ButtonWrapper),
	}
}

// AddButton 添加一个按钮
func (p *ButtonMgr) AddButton(id string, x, y, width, height int, btnText string, callback func() (ButtonAction, error)) {
	// 创建按钮
	btn := ebitenuiwidget.NewButton(
		ebitenuiwidget.ButtonOpts.WidgetOpts(
			ebitenuiwidget.WidgetOpts.MinSize(width, height),
		),
		ebitenuiwidget.ButtonOpts.Image(&ebitenuiwidget.ButtonImage{
			Idle:    ebitenuiimage.NewNineSliceColor(imgcolor.RGBA{R: 100, G: 100, B: 200, A: 255}),
			Hover:   ebitenuiimage.NewNineSliceColor(imgcolor.RGBA{R: 120, G: 120, B: 220, A: 255}),
			Pressed: ebitenuiimage.NewNineSliceColor(imgcolor.RGBA{R: 80, G: 80, B: 180, A: 255}),
		}),
		// 添加文字
		ebitenuiwidget.ButtonOpts.Text(btnText, apcfont.GFaceButton, &ebitenuiwidget.ButtonTextColor{
			Idle: imgcolor.White, // 白色文字
		}),
		ebitenuiwidget.ButtonOpts.ClickedHandler(
			func(args *ebitenuiwidget.ButtonClickedEventArgs) {
				if callback != nil {
					action, err := callback()
					if err != nil {
						log.Printf("button clicked, id:%v err:%v", id, err)
						return
					}
					if wrapper, ok := p.buttons[id]; ok {
						switch action {
						case ButtonActionNone:
							// 什么都不做
						case ButtonActionHide:
							wrapper.visible = false
							wrapper.button.GetWidget().Disabled = true
						case ButtonActionShow:
							wrapper.visible = true
							wrapper.button.GetWidget().Disabled = false
						case ButtonActionDestroy:
							wrapper.visible = false
							wrapper.button.GetWidget().Disabled = true
							delete(p.buttons, id)
							GUIMgr.RemoveButton(wrapper.button)
						}
					}
				}
			},
		),
	)
	// 手动设置按钮位置
	btn.GetWidget().Rect.Min.X = x
	btn.GetWidget().Rect.Min.Y = y
	btn.GetWidget().Rect.Max.X = x + width
	btn.GetWidget().Rect.Max.Y = y + height
	// 保存按钮引用
	p.buttons[id] = &ButtonWrapper{
		id:       id,
		button:   btn,
		visible:  true,
		callback: callback,
	}
	// 添加按钮到 GUI 管理器
	GUIMgr.AddButton(btn)
}

// Update 更新 UI
func (p *ButtonMgr) Update() {
}

// Draw 绘制 UI
func (p *ButtonMgr) Draw(screen *ebiten.Image) {
}
