package battle

import (
	"airplaneClient/src/common"
)

// Fire 开火
func (p *Plane) Fire() *Bullet {
	if !p.IsEnemy() { // 用户
		// 计算子弹初始位置 - 在飞机顶端中央
		bulletX := p.GetX() + p.GetImageWidth()/2
		bulletY := p.GetY()

		// 创建一颗向上移动的子弹 (id=1, level=1, speed=5)
		return NewBullet(1, 1, bulletX, bulletY, 1, common.BulletDirectionUp, p)
	}
	// 敌人
	return nil
}
