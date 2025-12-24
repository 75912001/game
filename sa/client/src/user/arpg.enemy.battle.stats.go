package user

// 计算敌人战斗数值 todo menglc 临时算法
type ArpgEnemyBattleStats struct {
	enemy *ArpgEnemy
}

func NewArpgEnemyBattleStats(enemy *ArpgEnemy) *ArpgEnemyBattleStats {
	return &ArpgEnemyBattleStats{
		enemy: enemy,
	}
}

// GetHpMax 获取血量上限
func (p *ArpgEnemyBattleStats) GetHpMax() uint32 {
	return p.enemy.GetCfg().Attributes.HP * (1 + p.enemy.Level)
}

// GetMpMax 获取魔法上限
func (p *ArpgEnemyBattleStats) GetMpMax() uint32 {
	return 100
}

// GetAttack 获取攻击力
func (p *ArpgEnemyBattleStats) GetAttack() uint32 {
	return p.enemy.GetCfg().Attributes.Attack * (1 + p.enemy.Level)
}

// GetDefense 获取防御力
func (p *ArpgEnemyBattleStats) GetDefense(strength, endurance, agility, stamina uint32) uint32 {
	return p.enemy.GetCfg().Attributes.Defense * (1 + p.enemy.Level)
}

// GetAgility 获取敏捷
func (p *ArpgEnemyBattleStats) GetAgility(agility uint32) uint32 {
	return p.enemy.GetCfg().Attributes.Agility * (1 + p.enemy.Level)
}

// GetCritRate 获取暴击率
func (p *ArpgEnemyBattleStats) GetCritRate() float32 {
	return p.enemy.GetCfg().Attributes.CritRate
}

// GetCounterRate 获取反击率
func (p *ArpgEnemyBattleStats) GetCounterRate() float32 {
	return p.enemy.GetCfg().Attributes.CounterRate
}

// GetDodgeRate 获取闪避率
func (p *ArpgEnemyBattleStats) GetDodgeRate() float32 {
	return p.enemy.GetCfg().Attributes.DodgeRate
}

// GetHitRate 获取命中率
func (p *ArpgEnemyBattleStats) GetHitRate() float32 {
	return p.enemy.GetCfg().Attributes.HitRate
}

// GetCritDamageBonus 获取暴击伤害加成
func (p *ArpgEnemyBattleStats) GetCritDamageBonus() float32 {
	return p.enemy.GetCfg().Attributes.CritDamageBonus
}

// GetStatusResist 获取异常状态抗性比率
func (p *ArpgEnemyBattleStats) GetStatusResist() float32 {
	return p.enemy.GetCfg().Attributes.StatusResist
}
