package battle

import (
	"airplaneClient/src/common"
	"airplaneClient/src/resources"
	"math"
)

type Bullet struct {
	*common.Object
	direction float64 // 方向(弧度), 0为向右, π/2为向下, π为向左, 3π/2为向上
	owner     *Plane  // 发射这颗子弹的飞机
	damage    uint32  // 伤害值
}

const BulletDirectionUp = 3.0 * math.Pi / 2.0 // 子弹-方向-上
const BulletDirectionDown = math.Pi / 2.0     // 子弹-方向-下

// NewBullet 创建一个新子弹
func NewBullet(id, level uint32, isBelongsToUser bool, x, y, speed, direction float64, owner *Plane) *Bullet {
	frames, imageWidth, imageHeight, err := resources.LoadBulletFrames(id, level, 1)
	if err != nil {
		panic(err)
	}
	return &Bullet{
		Object: common.NewObject(id, level,
			isBelongsToUser,
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
		damage:    1,
	}
}

// GetOwner 获取发射这颗子弹的飞机
func (b *Bullet) GetOwner() *Plane {
	return b.owner
}

// IsOutOfScreen 判断子弹是否飞出屏幕
func (b *Bullet) IsOutOfScreen(screenWidth, screenHeight float64) bool {
	return b.GetX() < -b.GetImageWidth() || screenWidth < b.GetX() ||
		b.GetY() < -b.GetImageHeight() || screenHeight < b.GetY()
}

func (b *Bullet) GetDamage() uint32 {
	return b.damage
}
