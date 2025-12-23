package arpg

import (
	"math"
	"math/rand"
	"saClient/src/cfg"

	xtime "github.com/75912001/xlib/time"
)

// EnemySpawnPoint 刷怪点
type EnemySpawnPoint struct {
	Object      *cfg.TiledObject // 关联的 Tiled 对象 (包含 ID, EnemyGroupID, EnemyGroup, PatrolRadius, RespawnSecond 等)
	TiledMapCfg *cfg.TiledMap    // 地图配置引用

	// 缓存的世界坐标 (从 Object.X, Y 计算)
	WX, WY float32
	TX, TY float32

	// 运行时状态
	AllDeadTime int64    // 全部死亡时的时间戳 (秒)
	Enemies     []*Enemy // 当前存活的怪物
}

// NewEnemySpawnPoint 从 TiledObject 创建刷怪点
func NewEnemySpawnPoint(tiledMapCfg *cfg.TiledMap, obj *cfg.TiledObject) *EnemySpawnPoint {
	// 将 Tiled 像素坐标转换为 World 坐标
	th := float32(tiledMapCfg.TileHeight)
	tx := obj.X/th - 0.5
	ty := obj.Y/th - 0.5
	wx, wy := tiledMapCfg.IsometricCT.T2W(tx, ty)

	enemySpawnPoint := &EnemySpawnPoint{
		Object:      obj,
		TiledMapCfg: tiledMapCfg,
		WX:          wx,
		WY:          wy,
		TX:          tx,
		TY:          ty,
		AllDeadTime: 0,
		Enemies:     make([]*Enemy, 0),
	}

	return enemySpawnPoint
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

		wx := p.WX + float32(math.Cos(angle))*distance
		wy := p.WY + float32(math.Sin(angle))*distance

		// 限制在地图边界内
		wx, wy = p.TiledMapCfg.ClampMapBoundaryWithW(wx, wy)

		// 检查位置是否阻挡 (Tile 阻挡)
		tx, ty := p.TiledMapCfg.IsometricCT.W2T(wx, wy)
		if p.TiledMapCfg.IsBlockedByTileWithTF(tx, ty) {
			continue
		}

		// 检查位置是否在 Object 阻挡区域内
		if _, blocked := p.TiledMapCfg.FindBlockedByObject(wx, wy); blocked {
			continue
		}

		return wx, wy
	}
}

// GetAllEnemies 获取所有怪物
func (p *EnemySpawnPoint) GetAllEnemies() []*Enemy {
	return p.Enemies
}
