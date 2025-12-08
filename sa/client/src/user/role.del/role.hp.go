package role_del

import "saClient/src/user/role"

func (p *role.Role) GetHp() uint32 {
	return p.hp
}

func (p *role.Role) SetHp(hp uint32) {
	p.hp = hp
}

func (p *role.Role) IsDie() bool {
	return p.GetHp() == 0
}

// TakeDamage 受到伤害
func (p *role.Role) TakeDamage(damage uint32) {
	if damage < p.GetHp() {
		p.SetHp(p.GetHp() - damage)
	} else {
		p.SetHp(0)
	}
}
