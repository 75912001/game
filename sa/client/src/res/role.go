package res

import (
	"saClient/src/common"
	"saClient/src/proto"
)

// Role 角色资源数据
// 包含角色的所有动作动画数据
type Role struct {
	ID common.AssetID // 角色ID

	// Actions 角色动作映射表
	// key: proto.RoleAction (Move, AttackAxe, Hurt, Die 等)
	// value: 该动作的所有方向动画数据
	// 性能: map查找 O(1)
	// 扩展性: 添加新动作无需修改此结构, 只需添加配置
	Actions map[proto.RoleAction]*RoleActionData
}
