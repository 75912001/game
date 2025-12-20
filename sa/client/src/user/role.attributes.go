package user

import (
	"saClient/src/common"
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

// GetValueInt 获取属性值-int
func (p *Role) GetValueInt(assetID proto.AssetIDRecord) int {
	v, ok := p.roleRecord.AssetIDRecordMap[uint32(assetID)]
	if !ok {
		return 0
	}
	return int(v)
}

// GetValueU32 获取属性值-uint32
func (p *Role) GetValueU32(assetID proto.AssetIDRecord) uint32 {
	v, ok := p.roleRecord.AssetIDRecordMap[uint32(assetID)]
	if !ok {
		return 0
	}
	return uint32(v)
}

// GetValueU64 获取属性值-uint64
func (p *Role) GetValueU64(assetID proto.AssetIDRecord) uint64 {
	v, ok := p.roleRecord.AssetIDRecordMap[uint32(assetID)]
	if !ok {
		return 0
	}
	return v
}

// SetValueU64 设置属性值
func (p *Role) SetValueU64(assetID proto.AssetIDRecord, value uint64) {
	p.roleRecord.AssetIDRecordMap[uint32(assetID)] = value
}

// GetValueF32 获取属性值-float32 (/1000倍)
func (p *Role) GetValueF32(assetID proto.AssetIDRecord) float32 {
	v, ok := p.roleRecord.AssetIDRecordMap[uint32(assetID)]
	if !ok {
		return 0
	}
	return float32(v) / float32(common.Float32Ratio)
}

// SetValueF32 设置属性值-float32 (*1000倍)
func (p *Role) SetValueF32(assetID proto.AssetIDRecord, value float32) {
	p.roleRecord.AssetIDRecordMap[uint32(assetID)] = uint64(value * float32(common.Float32Ratio))
}
