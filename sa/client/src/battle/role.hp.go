package battle

func (p *Role) GetHp() uint32 {
	return p.hp
}

func (p *Role) SetHp(hp uint32) {
	p.hp = hp
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
