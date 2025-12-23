package coordinatetransform

import (
	"math"
	"saClient/src/proto"
)

// CalculateOrientation 根据移动方向计算朝向
func CalculateOrientation(dx, dy float32) uint32 {
	// 8 方向判断
	angle := math.Atan2(float64(dy), float64(dx))
	// 转换为 0-360 度
	degrees := angle * 180 / math.Pi
	if degrees < 0 {
		degrees += 360
	}

	// 根据角度确定方向
	// 右: -22.5 ~ 22.5, 右下: 22.5 ~ 67.5, 下: 67.5 ~ 112.5, 左下: 112.5 ~ 157.5
	// 左: 157.5 ~ 202.5, 左上: 202.5 ~ 247.5, 上: 247.5 ~ 292.5, 右上: 292.5 ~ 337.5
	switch {
	case degrees < 22.5 || degrees >= 337.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Right)
	case degrees < 67.5:
		return uint32(proto.AssetOrientation_AssetOrientation_DownRight)
	case degrees < 112.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Down)
	case degrees < 157.5:
		return uint32(proto.AssetOrientation_AssetOrientation_DownLeft)
	case degrees < 202.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Left)
	case degrees < 247.5:
		return uint32(proto.AssetOrientation_AssetOrientation_UpLeft)
	case degrees < 292.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Up)
	default:
		return uint32(proto.AssetOrientation_AssetOrientation_UpRight)
	}
}
