package battle

import (
	"airplaneClient/src/resources"
	"math"
)

type Bullet struct {
	*Object
	direction float64 // 方向(弧度), 0为向右, π/2为向下, π为向左, 3π/2为向上
}

const BulletDirectionUp = 3.0 * math.Pi / 2.0
const BulletDirectionDown = math.Pi / 2.0

// NewBullet 创建一个新子弹
func NewBullet(id, level uint32, x, y, speed, direction float64) *Bullet {
	frames, imageWidth, imageHeight, err := resources.LoadPlaneFrames(id, level, 1)
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
	}
}
