package battle

import (
	"airplaneClient/src/common"
	"airplaneClient/src/resources"
	resourcescommon "airplaneClient/src/resources/common"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// Plane 飞机结构体
type Plane struct {
	*common.Object

	hp uint32

	currentFrameType resourcescommon.PlaneFrameType // 当前帧类型
	flipHorizontal   bool                           // 是否水平镜像(用于左倾)

	frameCounter int // 帧计数器，用于控制动画速度
}

// NewPlane 创建-飞机
func NewPlane(id, level uint32, isEnemy bool, x, y, speed, scale float64) *Plane {
	var frames []*ebitenv2.Image
	var imageWidth float64
	var imageHeight float64
	var err error
	var hp uint32
	var currentFrameType resourcescommon.PlaneFrameType
	if !isEnemy { // 用户
		// 加载飞机图片的4帧动画
		// 001.001.png 包含4帧:
		// 帧0-直飞(正), 帧1-微右倾, 帧2-右倾, 帧3-大右倾
		// 左倾通过镜像右倾帧实现
		var planeFrameType resourcescommon.PlaneFrameType
		frames, imageWidth, imageHeight, err = resources.LoadPlaneFrames(id, level, isEnemy, planeFrameType.GetFrameTypeCount())
		if err != nil {
			panic(err)
		}
		hp = 100
		currentFrameType = resourcescommon.PlaneFrameTypeStraight
	} else { // 敌人
		// 敌机只加载1帧图像
		frames, imageWidth, imageHeight, err = resources.LoadPlaneFrames(id, level, isEnemy, 1)
		if err != nil {
			panic(err)
		}
		hp = 5
		currentFrameType = resourcescommon.PlaneFrameTypeStraight
	}
	imageWidth *= scale
	imageHeight *= scale
	return &Plane{
		Object: common.NewObject(id, level,
			isEnemy,
			imageWidth,
			imageHeight,
			imageWidth*0.6,  // 碰撞体宽度为图像宽度的60%
			imageHeight*0.6, // 碰撞体高度为图像高度的60%
			x,
			y,
			speed,
			scale, // 缩放比例
			frames,
		),
		hp:               hp,
		currentFrameType: currentFrameType,
		flipHorizontal:   false, // 默认不镜像
		frameCounter:     0,
	}
}

// GetCurrentImage 获取当前飞机图像
func (p *Plane) GetCurrentImage() *ebitenv2.Image {
	return p.GetFrames()[p.currentFrameType]
}

// IsOutOfScreen 判断敌机是否飞出屏幕
func (p *Plane) IsOutOfScreen(screenHeight float64) bool {
	if !p.IsEnemy() { // 用户
		return false
	}
	return screenHeight < p.GetY()
}
