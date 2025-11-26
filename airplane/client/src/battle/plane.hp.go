package battle

func (p *Plane) GetHp() uint32 {
	return p.GetHp()
}

func (p *Plane) IsDestroyed() bool {
	return p.GetHp() == 0
}

// TakeDamage 飞机受到伤害
func (p *Plane) TakeDamage(damage uint32) {
	if damage < p.GetHp() {
		p.SetHp(p.GetHp() - damage)
	} else {
		p.SetHp(0)
	}
}
