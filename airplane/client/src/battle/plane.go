package battle

import (
	"airplaneClient/src/resources"
	resourcescommon "airplaneClient/src/resources/common"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// Plane 飞机结构体
type Plane struct {
	*Object

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
		Object: newObject(id, level,
			imageWidth,
			imageHeight,
			imageWidth*0.6,  // 碰撞体宽度为图像宽度的60%
			imageHeight*0.6, // 碰撞体高度为图像高度的60%
			x,
			y,
			speed,
			frames,
		),
		currentFrameType: resourcescommon.PlaneFrameTypeStraight,
		flipHorizontal:   false, // 默认不镜像
		frameCounter:     0,
	}
}

// GetCurrentImage 获取当前飞机图像
func (p *Plane) GetCurrentImage() *ebitenv2.Image {
	return p.frames[p.currentFrameType]
}
