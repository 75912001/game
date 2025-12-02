package common

type Object struct {
	image *Image

	x float64 // x 坐标
	y float64 // y 坐标

	scale float64 // 缩放比例(1.0为原始大小)

	// 是否画出图像边界(调试用) 蓝色
	debugDrawImageBounds bool
}
