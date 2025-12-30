package common

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

// 12 种颜色
var Colors_Red = color.RGBA{R: 255, G: 0, B: 0, A: 255}       // 红色
var Colors_Green = color.RGBA{R: 0, G: 255, B: 0, A: 255}     // 绿色
var Colors_Blue = color.RGBA{R: 0, G: 0, B: 255, A: 255}      // 蓝色
var Colors_Yellow = color.RGBA{R: 255, G: 255, B: 0, A: 255}  // 黄色
var Colors_Cyan = color.RGBA{R: 0, G: 255, B: 255, A: 255}    // 青色
var Colors_Magenta = color.RGBA{R: 255, G: 0, B: 255, A: 255} // 品红
var Colors_Orange = color.RGBA{R: 255, G: 165, B: 0, A: 255}  // 橙色
var Colors_Purple = color.RGBA{R: 128, G: 0, B: 128, A: 255}  // 紫色
var Colors_Pink = color.RGBA{R: 255, G: 192, B: 203, A: 255}  // 粉色
var Colors_Brown = color.RGBA{R: 165, G: 42, B: 42, A: 255}   // 棕色
var Colors_Gray = color.RGBA{R: 128, G: 128, B: 128, A: 255}  // 灰色
var Colors_Black = color.RGBA{R: 0, G: 0, B: 0, A: 255}       // 黑色

// DrawDashedCircle 绘制虚线圆
func DrawDashedCircle(screen *ebitenv2.Image, cx, cy, radius float32, clr color.RGBA, dashCount int, dashRatio float32, strokeWidth float32) {
	if dashCount <= 0 || radius <= 0 {
		return
	}

	angleStep := 2 * math.Pi / float64(dashCount)
	dashAngle := angleStep * float64(dashRatio)

	for i := 0; i < dashCount; i++ {
		startAngle := float64(i) * angleStep
		endAngle := startAngle + dashAngle

		// 计算起点和终点
		x1 := cx + radius*float32(math.Cos(startAngle))
		y1 := cy + radius*float32(math.Sin(startAngle))
		x2 := cx + radius*float32(math.Cos(endAngle))
		y2 := cy + radius*float32(math.Sin(endAngle))

		vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, clr, false)
	}
}
