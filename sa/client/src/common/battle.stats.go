package common

var GRoleBattleStats = &RoleBattleStats{}

type RoleBattleStats struct {
}

// GetHpMax 获取血量上限
func (p *RoleBattleStats) GetHpMax(strength, endurance, agility, stamina uint32) uint32 {
	return strength*1 + endurance*1 + agility*1 + stamina*4
}

// GetMpMax 获取魔法上限
func (p *RoleBattleStats) GetMpMax() uint32 {
	return 100
}

// GetAttack 获取攻击力
func (p *RoleBattleStats) GetAttack(strength, endurance, agility, stamina uint32) uint32 {
	// 使用整数运算: 乘以100后再除以100
	return uint32((strength*100 + endurance*10 + agility*5 + stamina*10) / 100)
}

// GetDefense 获取防御力
func (p *RoleBattleStats) GetDefense(strength, endurance, agility, stamina uint32) uint32 {
	// 使用整数运算: 乘以100后再除以100
	return uint32((strength*10 + endurance*100 + agility*5 + stamina*10) / 100)
}

// GetAgility 获取敏捷
func (p *RoleBattleStats) GetAgility(agility uint32) uint32 {
	return agility
}
