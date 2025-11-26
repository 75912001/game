package game

import (
	"airplaneClient/src/battle"
	"airplaneClient/src/common"
	"math/rand"
)

// enemySpawnCounter 敌机生成计数器
var enemySpawnCounter = 0

// updateEnemyPlanes 更新所有敌机
func (p *Game) updateEnemyPlanes() {
	// 更新每架敌机的位置
	for _, plane := range p.enemyPlanes {
		plane.Update()
	}

	// 移除飞出屏幕或被摧毁的敌机
	validPlanes := make([]*battle.Plane, 0, len(p.enemyPlanes))
	for _, plane := range p.enemyPlanes {
		if plane.IsDestroyed() {
			continue
		}
		if plane.IsOutOfScreen(common.ScreenHeight) {
			continue
		}
		validPlanes = append(validPlanes, plane)
	}
	p.enemyPlanes = validPlanes
}

// spawnEnemyPlane 生成敌机
func (p *Game) spawnEnemyPlane() {
	enemySpawnCounter++

	// 每240帧（约4秒）生成一架敌机
	if 240 <= enemySpawnCounter {
		enemySpawnCounter = 0

		// 在屏幕顶部随机位置生成敌机
		x := float64(rand.Intn(int(common.ScreenWidth - 50)))
		y := -50.0 // 从屏幕上方生成

		// 创建敌机 (id=1, level=1, speed=1缓慢向下, scale=0.8小一点)
		enemy := battle.NewPlane(1, 1, true, x, y, 1, 0.2)
		p.enemyPlanes = append(p.enemyPlanes, enemy)
	}
}
