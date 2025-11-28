package resources

import (
	"airplaneClient/src/common"
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"path/filepath"
)

// GenPlaneName 生成-飞机-资源名称
func GenPlaneName(id uint32, level uint32) string {
	return fmt.Sprintf("plane.%03d.%03d.png", id, level)
}

// LoadPlaneFrames 加载-飞机-动画帧
// 从一张包含多帧的图片中切分出每一帧
func LoadPlaneFrames(id uint32, level uint32, isEnemy bool, frameCount uint32) (frames []*ebitenv2.Image, frameWidth float64, frameHeight float64, err error) {
	// 构建图片路径
	var planeName string
	if !isEnemy { // 用户
		planeName = GenPlaneName(id, level)
	} else { // 敌人
		planeName = GenEnemyPlaneName(id, level)
	}
	planePath := filepath.Join(common.AppResourcesDir, "planes", planeName)

	// 加载图片
	img, _, err := ebitenv2ebitenutil.NewImageFromFile(planePath)
	if err != nil {
		log.Printf("加载飞机图片失败:%v %v", planePath, err)
		return nil, 0, 0, err
	}

	// 获取图片尺寸
	bounds := img.Bounds()
	totalWidth := bounds.Dx()
	totalHeight := bounds.Dy()

	// 计算每一帧的宽度
	frameWidth = float64(totalWidth) / float64(frameCount)
	frameHeight = float64(totalHeight)

	// 切分帧
	frames = make([]*ebitenv2.Image, frameCount)
	for i := 0; i < int(frameCount); i++ {
		// 创建子图像区域
		subImg := img.SubImage(image.Rect(
			i*int(frameWidth), 0,
			(i+1)*int(frameWidth), int(frameHeight),
		)).(*ebitenv2.Image)

		frames[i] = subImg
	}

	return frames, frameWidth, frameHeight, nil
}
