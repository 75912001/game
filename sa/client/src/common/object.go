package common

import (
	"saClient/src/proto"
)

type Object struct {
	image    *Image
	frameIdx uint32 // 当前帧索引

	Point proto.Point

	debugDrawImageBounds bool // 是否画出图像边界(调试用)
}
