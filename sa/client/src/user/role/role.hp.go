package role

import (
	"saClient/src/proto"
)

// GetHp 获取血量
func (p *Role) GetHp() uint32 {
	v, _ := p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)]
	return uint32(v)
}

// SetHp 设置血量
func (p *Role) SetHp(hp uint32) {
	p.roleRecord.AssetIDRecordMap[uint32(proto.AssetIDRecord_AssetIDRecord_HP)] = uint64(hp)
}

func (p *Role) IsDie() bool {
	return p.GetHp() == 0
}

// TakeDamage 受到伤害
func (p *Role) TakeDamage(damage uint32) {
	if damage < p.GetHp() {
		p.SetHp(p.GetHp() - damage)
	} else {
		p.SetHp(0)
	}
}
