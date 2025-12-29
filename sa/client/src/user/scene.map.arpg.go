package user

func (p *Map) IsArpgMap() bool {
	return 0 < len(p.spawnManager.spawnPoints)
}
