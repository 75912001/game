package arpg

import (
	"saClient/src/cfg"
)

// EnemySpawnMgr 敌人-刷怪管理器
type EnemySpawnMgr struct {
	spawnPoints []*EnemySpawnPoint // 刷怪点列表
	tiledMapCfg *cfg.TiledMap      // 地图配置引用
}

// NewEnemySpawnMgr 创建刷怪管理器
func NewEnemySpawnMgr(tiledMapCfg *cfg.TiledMap) *EnemySpawnMgr {
	sm := &EnemySpawnMgr{
		spawnPoints: make([]*EnemySpawnPoint, 0),
		tiledMapCfg: tiledMapCfg,
	}

	// 从地图中查找所有刷怪点
	spawnPointObjs := tiledMapCfg.GetSpawnEnemyGroupPoints()
	for _, obj := range spawnPointObjs {
		sp := NewEnemySpawnPoint(tiledMapCfg, obj)
		sm.spawnPoints = append(sm.spawnPoints, sp)
		sp.SpawnEnemy()
	}

	return sm
}

// Update 每帧更新
func (p *EnemySpawnMgr) Update() {
	for _, sp := range p.spawnPoints {
		sp.Update()
	}
}

// GetAllEnemies 获取所有怪物 (用于渲染)
func (p *EnemySpawnMgr) GetAllEnemies() []*Enemy {
	var result []*Enemy
	for _, sp := range p.spawnPoints {
		result = append(result, sp.GetAllEnemies()...)
	}
	return result
}
