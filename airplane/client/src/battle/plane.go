package battle

import (
	"airplaneClient/src/resources"
	resourcescommon "airplaneClient/src/resources/common"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// Plane 飞机结构体
type Plane struct {
	id    uint32 // 飞机ID
	level uint32 // 飞机等级

	imageWidth     float64 // 图像-宽
	imageHeight    float64 // 图像-高
	colliderWidth  float64 // 碰撞体-宽
	colliderHeight float64 // 碰撞体-高

	x     float64 // x 坐标
	y     float64 // y 坐标
	speed float64 // 移动速度

	frames           []*ebitenv2.Image              // 动画帧：0-直飞, 1-微右倾, 2-右倾, 3-大右倾
	currentFrameType resourcescommon.PlaneFrameType // 当前帧类型
	flipHorizontal   bool                           // 是否水平镜像(用于左倾)

	frameCounter int // 帧计数器，用于控制动画速度
}

// NewPlane 创建一个新飞机
func NewPlane(id, level uint32, x, y, speed float64) *Plane {
	// 加载飞机图片的4帧动画
	// 001.001.png 包含4帧:
	// 帧0-直飞(正), 帧1-微右倾, 帧2-右倾, 帧3-大右倾
	// 左倾通过镜像右倾帧实现
	var planeFrameType resourcescommon.PlaneFrameType
	frames, imageWidth, imageHeight, err := resources.LoadPlaneFrames(id, level, planeFrameType.GetFrameTypeCount())
	if err != nil {
		panic(err)
	}
	return &Plane{
		id:    id,
		level: level,

		imageWidth:     imageWidth,
		imageHeight:    imageHeight,
		colliderWidth:  imageWidth * 0.6,  // 碰撞体宽度为图像宽度的60%
		colliderHeight: imageHeight * 0.6, // 碰撞体高度为图像高度的60%

		x:     x,
		y:     y,
		speed: speed,

		frames:           frames,
		currentFrameType: resourcescommon.PlaneFrameTypeStraight,
		flipHorizontal:   false, // 默认不镜像

		frameCounter: 0,
	}
}

// GetCurrentImage 获取当前飞机图像
func (p *Plane) GetCurrentImage() *ebitenv2.Image {
	return p.frames[p.currentFrameType]
}
