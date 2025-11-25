package battle

import ebitenv2 "github.com/hajimehoshi/ebiten/v2"

type Object struct {
	id    uint32 // ID
	level uint32 // 等级

	imageWidth  float64 // 图像-宽
	imageHeight float64 // 图像-高

	colliderWidth  float64 // 碰撞体-宽
	colliderHeight float64 // 碰撞体-高

	x float64 // x 坐标
	y float64 // y 坐标

	speed float64 // 移动速度

	frames []*ebitenv2.Image // 动画帧
}

func newObject(id, level uint32,
	imageWidth, imageHeight, colliderWidth, colliderHeight float64,
	x, y, speed float64, frames []*ebitenv2.Image) *Object {
	return &Object{
		id:    id,
		level: level,

		imageWidth:  imageWidth,
		imageHeight: imageHeight,

		colliderWidth:  colliderWidth,
		colliderHeight: colliderHeight,

		x: x,
		y: y,

		speed: speed,

		frames: frames,
	}
}
