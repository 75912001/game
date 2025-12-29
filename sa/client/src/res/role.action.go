package res

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

// RoleActionData 角色-动作数据 (通用结构,适用于所有动作类型)
// 用于存储单个动作的所有方向的动画帧数据
// 性能: 数组访问 O(1), 内存布局紧凑
type RoleActionData struct {
	// Frames 存储每个方向的动画帧图像
	// 索引: proto.AssetOrientation 枚举值
	// 镜像方向(右/右上/右下)的帧由左侧方向自动生成
	Frames [proto.AssetOrientation_AssetOrientation_Max][]*ebitenv2.Image

	// FrameInfo 存储每个方向的动画帧配置信息
	// 与 Frames 一一对应, 包含帧的尺寸等元数据
	// 镜像帧信息与原帧相同(因为只是水平翻转)
	FrameInfo [proto.AssetOrientation_AssetOrientation_Max][]*RoleImageSprite
}

// NewRoleActionData 创建新的角色动作数据
// 返回空的动作数据结构, 等待资源加载填充
func NewRoleActionData() *RoleActionData {
	return &RoleActionData{}
}
