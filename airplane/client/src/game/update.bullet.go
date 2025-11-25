package game

import (
	"airplaneClient/src/battle"
	"airplaneClient/src/common"
)

// updateBullets 更新所有子弹
func (p *Game) updateBullets() {
	{ // 更新用户子弹
		for _, bullet := range p.userBullets {
			bullet.Update()
		}
		// 子弹与飞机的碰撞
		p.userBullets = p.bulletCollisions(p.userBullets, p.enemyPlanes)
		// 移除飞出屏幕的子弹
		validBullets := make([]*battle.Bullet, 0, len(p.userBullets))
		for _, bullet := range p.userBullets {
			if !bullet.IsOutOfScreen(common.ScreenWidth, common.ScreenHeight) {
				validBullets = append(validBullets, bullet)
			}
		}
		p.userBullets = validBullets
	}
	{ // 更新敌人子弹
		for _, bullet := range p.enemyBullets {
			bullet.Update()
		}
		// 弹与飞机的碰撞
		p.enemyBullets = p.bulletCollisions(p.enemyBullets, []*battle.Plane{p.userPlane})
		// 移除飞出屏幕的子弹
		validBullets := make([]*battle.Bullet, 0, len(p.enemyBullets))
		for _, bullet := range p.enemyBullets {
			if !bullet.IsOutOfScreen(common.ScreenWidth, common.ScreenHeight) {
				validBullets = append(validBullets, bullet)
			}
		}
		p.enemyBullets = validBullets
	}
}

// 子弹与飞机的碰撞
func (p *Game) bulletCollisions(bullets []*battle.Bullet, planes []*battle.Plane) (validBullets []*battle.Bullet) {
	// 遍历所有子弹，检测是否击中飞机
	validBullets = make([]*battle.Bullet, 0, len(bullets))
	for _, bullet := range bullets {
		hit := false
		for _, plane := range planes {
			if plane.CollidesWith(bullet.Object) {
				// 子弹击中了飞机
				p.onPlaneHit(plane, bullet)
				hit = true
				break // 一个子弹只能击中一个目标，跳出循环
			}
		}
		if !hit { // 没有击中任何目标，保留子弹
			validBullets = append(validBullets, bullet)
		}
	}
	return validBullets
}

// onPlaneHit 当飞机被子弹击中时调用
func (p *Game) onPlaneHit(plane *battle.Plane, bullet *battle.Bullet) {
	// 飞机受到伤害
	plane.TakeDamage(bullet.GetDamage())

	// 检查飞机是否被摧毁
	if plane.IsDestroyed() { // 飞机被摧毁
		if plane == p.userPlane { // 用户飞机
			p.state = StateGameOver
		}
	}
}
