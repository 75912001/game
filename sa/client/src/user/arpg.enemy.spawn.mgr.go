package user

import (
	"saClient/src/cfg"
)

// ArpgEnemySpawnMgr 敌人-刷怪管理器
type ArpgEnemySpawnMgr struct {
	spawnPoints []*ArpgEnemySpawnPoint // 刷怪点列表
	tiledMapCfg *cfg.TiledMap          // 地图配置引用
}

// NewArpgEnemySpawnMgr 创建刷怪管理器
func NewArpgEnemySpawnMgr(tiledMapCfg *cfg.TiledMap) *ArpgEnemySpawnMgr {
	sm := &ArpgEnemySpawnMgr{
		spawnPoints: make([]*ArpgEnemySpawnPoint, 0),
		tiledMapCfg: tiledMapCfg,
	}

	// 从地图中查找所有刷怪点
	spawnPointObjs := tiledMapCfg.GetSpawnEnemyGroupPoints()
	for _, obj := range spawnPointObjs {
		sp := NewArpgEnemySpawnPoint(tiledMapCfg, obj)
		sm.spawnPoints = append(sm.spawnPoints, sp)
		sp.SpawnEnemy()
	}

	return sm
}

// Update 每帧更新
func (p *ArpgEnemySpawnMgr) Update() {
	for _, sp := range p.spawnPoints {
		sp.Update()
	}
}

// GetAllEnemies 获取所有怪物 (用于渲染)
func (p *ArpgEnemySpawnMgr) GetAllEnemies() []*ArpgEnemy {
	var result []*ArpgEnemy
	for _, sp := range p.spawnPoints {
		result = append(result, sp.GetAllEnemies()...)
	}
	return result
}
