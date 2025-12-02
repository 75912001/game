package battle

import (
	"saClient/src/common"
	"saClient/src/proto"
)

type User struct {
	object *common.Object
	record *proto.RoleRecord // 角色-记录

	speed float64 // 移动速度
	hp    uint32  // 生命值
	scale float64 // 缩放比例(1.0为原始大小)
}
