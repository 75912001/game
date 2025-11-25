package battle

import (
	"airplaneClient/src/resources"
)

type Bullet struct {
	*Object
}

// NewBullet 创建一个新子弹
func NewBullet(id, level uint32, x, y, speed float64) *Bullet {
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
	}
}
