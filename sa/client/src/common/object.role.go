package common

const RoleSpeed float64 = 2.0 // 角色默认移动速度

type ObjectRole struct {
	id uint32 // ID

	object *Object

	elemental *Elemental // 元素属性

	speed float64 // 移动速度
	hp    uint32  // 生命值

	attribute *RoleAttribute
}

// RoleAttribute 角色属性
type RoleAttribute struct {
	level   uint32 // 等级
	exp     uint64 // 经验值
	attack  uint32 // 攻
	defense uint32 // 防
	agility uint32 // 敏
	hp      uint32 // 生命
}
