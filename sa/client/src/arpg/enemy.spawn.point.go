package arpg

import (
	"math"
	"math/rand"
	"saClient/src/cfg"

	xtime "github.com/75912001/xlib/time"
)

// EnemySpawnPoint 刷怪点
type EnemySpawnPoint struct {
	Object      *cfg.TiledObject // 关联的 Tiled 对象
	TiledMapCfg *cfg.TiledMap    // 地图配置引用

	// 运行时状态
	AllDeadTime int64    // 全部死亡时的时间戳 (秒)
	Enemies     []*Enemy // 当前存活的怪物
}

// NewEnemySpawnPoint 从 TiledObject 创建刷怪点
func NewEnemySpawnPoint(tiledMapCfg *cfg.TiledMap, obj *cfg.TiledObject) *EnemySpawnPoint {
	return &EnemySpawnPoint{
		Object:      obj,
		TiledMapCfg: tiledMapCfg,
		AllDeadTime: 0,
		Enemies:     make([]*Enemy, 0),
	}
}

// Update 每帧更新
func (p *EnemySpawnPoint) Update() {
	// 1. 更新所有怪物 AI
	for _, enemy := range p.Enemies {
		enemy.Update()
	}

	// 清理死亡怪物
	hadEnemies := 0 < len(p.Enemies)
	p.removeDeadEnemies()

	// 检查是否全部死亡，记录死亡时间
	if hadEnemies && len(p.Enemies) == 0 {
		p.AllDeadTime = xtime.GTimeMgr.ShadowTimestamp()
	}

	// 检查是否需要刷怪
	// RespawnSecond=0 表示不刷新
	if 0 < p.Object.RespawnSecond && len(p.Enemies) == 0 {
		if int64(p.Object.RespawnSecond) <= xtime.GTimeMgr.ShadowTimestamp()-p.AllDeadTime {
			p.SpawnEnemy()
			p.AllDeadTime = 0 // 重置
		}
	}
}

// removeDeadEnemies 清理死亡怪物
func (p *EnemySpawnPoint) removeDeadEnemies() {
	alive := p.Enemies[:0]
	for _, enemy := range p.Enemies {
		if !enemy.IsDead() {
			alive = append(alive, enemy)
		}
	}
	p.Enemies = alive
}

// SpawnEnemy 生成怪物
func (p *EnemySpawnPoint) SpawnEnemy() {
	// 根据怪物组配置生成敌人 (一次性生成全部)
	generated := p.Object.EnemyGroup.Generate(1) // todo menglc 传入玩家等级
	if len(generated) == 0 {
		return
	}

	for _, gen := range generated {
		// 在刷怪点附近随机位置生成
		wx, wy := p.randomSpawnPosition()
		enemy := NewEnemy(p, gen, wx, wy)
		p.Enemies = append(p.Enemies, enemy)
	}
}

// randomSpawnPosition 随机生成位置 (使用 SpawnRadius)
func (p *EnemySpawnPoint) randomSpawnPosition() (float32, float32) {
	return p.RandomPositionInRadius(p.Object.SpawnRadius)
}

// RandomPositionInRadius 在指定半径内随机生成一个位置(无阻挡,不一定可达)
func (p *EnemySpawnPoint) RandomPositionInRadius(radius float32) (float32, float32) {
	for {
		angle := rand.Float64() * 2 * math.Pi
		distance := rand.Float32() * radius

		wx := p.Object.WX + float32(math.Cos(angle))*distance
		wy := p.Object.WY + float32(math.Sin(angle))*distance

		// 综合检测 (地图边界 + Tile阻挡 + Object阻挡)
		wx, wy, blocked := p.TiledMapCfg.IsBlocked(wx, wy)
		if blocked {
			continue
		}

		return wx, wy
	}
}

// GetAllEnemies 获取所有怪物
func (p *EnemySpawnPoint) GetAllEnemies() []*Enemy {
	return p.Enemies
}
