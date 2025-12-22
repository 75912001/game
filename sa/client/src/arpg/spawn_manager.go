package arpg

import (
	"saClient/src/cfg"
)

// SpawnManager 刷怪点管理器
type SpawnManager struct {
	spawnPoints []*EnemyGroupSpawnPoint // 刷怪点列表
	tiledMapCfg *cfg.TiledMap           // 地图配置引用
}

// NewSpawnManager 创建刷怪点管理器
func NewSpawnManager(tiledMapCfg *cfg.TiledMap) *SpawnManager {
	sm := &SpawnManager{
		spawnPoints: make([]*EnemyGroupSpawnPoint, 0),
		tiledMapCfg: tiledMapCfg,
	}

	// 从地图中查找所有刷怪点
	spawnPointObjs := tiledMapCfg.FindSpawnEnemyGroupPoints()
	for _, obj := range spawnPointObjs {
		sp := NewEnemyGroupSpawnPoint(tiledMapCfg, obj)
		sm.spawnPoints = append(sm.spawnPoints, sp)
		// 初始生成怪物
		sp.SpawnEnemy()
	}

	return sm
}

// Update 每帧更新
func (p *SpawnManager) Update() {
	for _, sp := range p.spawnPoints {
		sp.Update()
	}
}

// GetAllEnemies 获取所有怪物 (用于渲染)
func (p *SpawnManager) GetAllEnemies() []*Enemy {
	var result []*Enemy
	for _, sp := range p.spawnPoints {
		result = append(result, sp.GetAllEnemies()...)
	}
	return result
}

// GetSpawnPointCount 获取刷怪点数量
func (p *SpawnManager) GetSpawnPointCount() int {
	return len(p.spawnPoints)
}

// GetEnemyCount 获取怪物总数
func (p *SpawnManager) GetEnemyCount() int {
	count := 0
	for _, sp := range p.spawnPoints {
		count += len(sp.Enemies)
	}
	return count
}
