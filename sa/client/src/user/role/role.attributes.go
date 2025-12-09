package role

import (
	"saClient/src/proto"
)

// GetCurrentHp 获取当前血量
func (p *Role) GetCurrentHp() uint32 {
	v, _ := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)]
	return uint32(v)
}

// SetCurrentHp 设置当前血量
func (p *Role) SetCurrentHp(hp uint32) {
	p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = uint64(hp)
}

func (p *Role) IsDie() bool {
	return p.GetCurrentHp() == 0
}

// TakeDamage 受到伤害
func (p *Role) TakeDamage(damage uint32) {
	if damage < p.GetCurrentHp() {
		p.SetCurrentHp(p.GetCurrentHp() - damage)
	} else {
		p.SetCurrentHp(0)
	}
}

// GetHpMax 获取血量上限
func (p *Role) GetHpMax() uint32 {
	strength := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Strength)]
	endurance := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Endurance)]
	agility := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Agility)]
	stamina := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Stamina)]
	return uint32(strength*1 + endurance*1 + agility*1 + stamina*4)
}

// GetAttack 获取攻击力
func (p *Role) GetAttack() uint32 {
	strength := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Strength)]
	endurance := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Endurance)]
	agility := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Agility)]
	stamina := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Stamina)]
	// 使用整数运算: 乘以100后再除以100
	return uint32((strength*100 + endurance*10 + agility*5 + stamina*10) / 100)
}

// GetDefense 获取防御力
func (p *Role) GetDefense() uint32 {
	strength := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Strength)]
	endurance := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Endurance)]
	agility := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Agility)]
	stamina := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Stamina)]
	// 使用整数运算: 乘以100后再除以100
	return uint32((strength*10 + endurance*100 + agility*5 + stamina*10) / 100)
}

// GetSpeed 获取速度
func (p *Role) GetSpeed() uint32 {
	agility := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_Agility)]
	return uint32(agility)
}
