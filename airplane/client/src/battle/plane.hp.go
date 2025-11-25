package battle

func (p *Plane) GetHp() uint32 {
	return p.hp
}

func (p *Plane) IsDestroyed() bool {
	return p.hp == 0
}

// TakeDamage 飞机受到伤害
func (p *Plane) TakeDamage(damage uint32) {
	if damage < p.hp {
		p.hp -= damage
	} else {
		p.hp = 0
	}
}
