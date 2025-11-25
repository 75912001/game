package battle

import (
	"airplaneClient/src/resources"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Bullet struct {
	*Object
	direction float64 // 方向(弧度), 0为向右, π/2为向下, π为向左, 3π/2为向上
	owner     *Plane  // 发射这颗子弹的飞机
}

const BulletDirectionUp = 3.0 * math.Pi / 2.0 // 子弹-方向-上
const BulletDirectionDown = math.Pi / 2.0     // 子弹-方向-下

// NewBullet 创建一个新子弹
func NewBullet(id, level uint32, x, y, speed, direction float64, owner *Plane) *Bullet {
	frames, imageWidth, imageHeight, err := resources.LoadBulletFrames(id, level, 1)
	if err != nil {
		panic(err)
	}
	return &Bullet{
		Object: newObject(id, level,
			imageWidth,
			imageHeight,
			imageWidth*0.3,  // 碰撞体宽度为图像宽度的30%
			imageHeight*0.3, // 碰撞体高度为图像高度的30%
			x,
			y,
			speed,
			frames,
		),
		direction: direction,
		owner:     owner,
	}
}

// GetOwner 获取发射这颗子弹的飞机
func (b *Bullet) GetOwner() *Plane {
	return b.owner
}

// Update 更新子弹位置
func (b *Bullet) Update() {
	// 根据方向更新位置
	b.x += math.Cos(b.direction) * b.speed
	b.y += math.Sin(b.direction) * b.speed
}

// Draw 绘制子弹
func (b *Bullet) Draw(screen *ebitenv2.Image) {
	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.frames[0], op)
}

// IsOutOfScreen 判断子弹是否飞出屏幕
func (b *Bullet) IsOutOfScreen(screenWidth, screenHeight float64) bool {
	return b.x < -b.imageWidth || screenWidth < b.x ||
		b.y < -b.imageHeight || screenHeight < b.y
}
