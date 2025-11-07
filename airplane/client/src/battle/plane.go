package battle

import (
	apccommon "airplaneClient/src/common"
	apcui "airplaneClient/src/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

var planeWidth, planeHeight = 10, 30

// Plane 飞机结构体
type Plane struct {
	x     float64 // X 坐标
	y     float64 // Y 坐标
	speed float64 // 移动速度
	image *ebiten.Image
}

// NewPlane 创建一个新飞机
func NewPlane(x, y, speed float64) *Plane {
	// todo menglc 加载实际的飞机图片

	planeImg := ebiten.NewImage(planeWidth, planeHeight)
	planeImg.Fill(color.RGBA{R: 255, G: 255, B: 0, A: 255}) // 黄色飞机

	return &Plane{
		x:     x,
		y:     y,
		speed: speed,
		image: planeImg,
	}
}

// Update 更新飞机状态
func (p *Plane) Update() {
	// 键盘控制
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.x -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.x += p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		p.y -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		p.y += p.speed
	}

	// 边界检测
	if p.x < 0 {
		p.x = 0
	}
	if p.x > float64(apccommon.ScreenWidth-planeWidth) {
		p.x = float64(apccommon.ScreenWidth - planeWidth)
	}
	if p.y < 0 {
		p.y = 0
	}
	if p.y > float64(apccommon.ScreenHeight-planeHeight) {
		p.y = float64(apccommon.ScreenHeight - planeHeight)
	}
}

// Draw 绘制飞机
func (p *Plane) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.image, op)

	apcui.Printf(screen, 10, 10, "使用方向键或WASD移动飞机")
}
