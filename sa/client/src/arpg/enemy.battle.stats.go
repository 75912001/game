package arpg

// 计算敌人战斗数值
type EnemyBattleStats struct {
	enemy *Enemy
}

func NewEnemyBattleStats(enemy *Enemy) *EnemyBattleStats {
	return &EnemyBattleStats{
		enemy: enemy,
	}
}

// GetHpMax 获取血量上限
func (p *EnemyBattleStats) GetHpMax() uint32 {
	return p.enemy.GetCfg().Attributes.HP * (1 + p.enemy.Level)
}

// GetMpMax 获取魔法上限
func (p *EnemyBattleStats) GetMpMax() uint32 {
	return 100
}

// GetAttack 获取攻击力
func (p *EnemyBattleStats) GetAttack() uint32 {
	return p.enemy.GetCfg().Attributes.Attack * (1 + p.enemy.Level)
}

// GetDefense 获取防御力
func (p *EnemyBattleStats) GetDefense(strength, endurance, agility, stamina uint32) uint32 {
	return p.enemy.GetCfg().Attributes.Defense * (1 + p.enemy.Level)
}

// GetAgility 获取敏捷
func (p *EnemyBattleStats) GetAgility(agility uint32) uint32 {
	return p.enemy.GetCfg().Attributes.Agility * (1 + p.enemy.Level)
}

// GetCritRate 获取暴击率
func (p *EnemyBattleStats) GetCritRate() float32 {
	return p.enemy.GetCfg().Attributes.CritRate
}

// GetCounterRate 获取反击率
func (p *EnemyBattleStats) GetCounterRate() float32 {
	return p.enemy.GetCfg().Attributes.CounterRate
}

// GetDodgeRate 获取闪避率
func (p *EnemyBattleStats) GetDodgeRate() float32 {
	return p.enemy.GetCfg().Attributes.DodgeRate
}

// GetHitRate 获取命中率
func (p *EnemyBattleStats) GetHitRate() float32 {
	return p.enemy.GetCfg().Attributes.HitRate
}

// GetCritDamageBonus 获取暴击伤害加成
func (p *EnemyBattleStats) GetCritDamageBonus() float32 {
	return p.enemy.GetCfg().Attributes.CritDamageBonus
}

// GetStatusResist 获取异常状态抗性比率
func (p *EnemyBattleStats) GetStatusResist() float32 {
	return p.enemy.GetCfg().Attributes.StatusResist
}
