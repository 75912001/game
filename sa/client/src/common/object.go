package common

import (
	"saClient/src/proto"
)

type Object struct {
	image    *Image
	frameIdx uint32 // 当前帧索引

	Point proto.Point

	// 是否画出图像边界(调试用) 蓝色
	debugDrawImageBounds bool
}
