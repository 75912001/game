package user

import (
	"saClient/src/proto"
)

// 计算角色战斗数值
type RoleBattleStats struct {
	role *Role
}

func NewRoleBattleStats(role *Role) *RoleBattleStats {
	return &RoleBattleStats{
		role: role,
	}
}

// GetHpMax 获取血量上限
func (p *RoleBattleStats) GetHpMax() uint32 {
	strength := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStrength)])
	endurance := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesEndurance)])
	agility := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesAgility)])
	stamina := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStamina)])
	return strength*1 + endurance*1 + agility*1 + stamina*4
}

// GetMpMax 获取魔法上限
func (p *RoleBattleStats) GetMpMax() uint32 {
	return 100
}

// GetAttack 获取攻击力
func (p *RoleBattleStats) GetAttack() uint32 {
	strength := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStrength)])
	endurance := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesEndurance)])
	agility := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesAgility)])
	stamina := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStamina)])
	// 使用整数运算: 乘以100后再除以100
	return (strength*100 + endurance*10 + agility*5 + stamina*10) / 100
}

// GetDefense 获取防御力
func (p *RoleBattleStats) GetDefense() uint32 {
	strength := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStrength)])
	endurance := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesEndurance)])
	agility := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesAgility)])
	stamina := uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesStamina)])

	// 使用整数运算: 乘以100后再除以100
	return (strength*10 + endurance*100 + agility*5 + stamina*10) / 100
}

// GetAgility 获取敏捷
func (p *RoleBattleStats) GetAgility() uint32 {
	return uint32(p.role.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Role_AttributesAgility)])
}

// GetCritRate 获取暴击率
func (p *RoleBattleStats) GetCritRate() float32 {
	return 0.12
}

// GetCounterRate 获取反击率
func (p *RoleBattleStats) GetCounterRate() float32 {
	return 0.08
}

// GetDodgeRate 获取闪避率
func (p *RoleBattleStats) GetDodgeRate() float32 {
	return 0.15
}

// GetHitRate 获取命中率
func (p *RoleBattleStats) GetHitRate() float32 {
	return 0.92
}

// GetCritDamageBonus 获取暴击伤害加成
func (p *RoleBattleStats) GetCritDamageBonus() float32 {
	return 0.6
}

// GetStatusResist 获取异常状态抗性比率
func (p *RoleBattleStats) GetStatusResist() float32 {
	return 0.25
}
