package common

type Object struct {
	image *Image

	framesIdx int // 当前帧索引

	x float64 // x 坐标
	y float64 // y 坐标

	// 是否画出图像边界(调试用) 蓝色
	debugDrawImageBounds bool
}
