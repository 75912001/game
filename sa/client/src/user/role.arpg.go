package user

import (
	xtime "github.com/75912001/xlib/time"
	"math"
	"saClient/src/cfg"
	commonct "saClient/src/common/coordinatetransform"
	"saClient/src/proto"
)

type RoleArpg struct {
	TargetEnemy         *ArpgEnemy // 攻击目标
	AttackCDMillisecond int64      // 攻击冷却计时器(毫秒)
}

func NewRoleArpg() *RoleArpg {
	return &RoleArpg{}
}

// handleArpgAutoAttack 处理自动攻击逻辑
func (p *Role) handleArpgAutoAttack() {
	// 1. 检查当前状态
	// 如果正在攻击且动画未播放完，保持攻击状态 (锁定)
	if p.sprite.action == proto.RoleAction_RoleAction_AttackAxe {
		p.animationFrame.Update()

		actionData := p.cfgRole.ResRole.Actions[p.sprite.action]
		frames := actionData.Frames[p.sprite.orientation]
		// 简单判断：如果帧索引超过了总帧数，说明播放了一轮
		if p.animationFrame.FrameIdx >= uint32(len(frames))-1 { // 检查动画播放完毕
			p.SetAction(proto.RoleAction_RoleAction_Move) // 攻击结束，切回默认状态
			p.UpdateWithAction()                          // 更新图像
		} else { // 继续播放攻击动画
			p.UpdateWithAction() // 更新图像
			return
		}
	}

	p.selectTarget()
	if p.Arpg.TargetEnemy == nil { // 没有目标
		return
	}
	// 检查距离
	dist := p.distanceToEnemy(p.Arpg.TargetEnemy)

	attackRange := cfg.GCommon.GetRoleArpgDefAttackRangeByWeaponType(proto.RoleWeaponType_RoleWeaponType_Axe)
	if dist <= attackRange {
		// 检查冷却
		nowMillisecond := xtime.GTimeMgr.GetMillisecond()
		if nowMillisecond >= p.Arpg.AttackCDMillisecond {
			p.performAttack()
			p.Arpg.AttackCDMillisecond = nowMillisecond + 1000 // 假设攻击间隔 1秒
			p.UpdateWithAction()                               // 更新图像
		}
	} else {
		// 目标在范围外，这里可以添加自动追击逻辑，或者什么都不做保持原地
	}
}

// selectTarget 选择攻击目标
func (p *Role) selectTarget() {
	if p.Arpg.TargetEnemy != nil { // 已有目标
		if p.Arpg.TargetEnemy.IsDead() { // 目标死亡
			p.Arpg.TargetEnemy = nil
		} else { // 目标存活
			return
		}
	}
	// 寻找最近的敌人
	var nearestEnemy *ArpgEnemy
	var minDist float32 = math.MaxFloat32 // 最小距离
	searchRange := cfg.GCommon.RoleArpgDefViewRange

	// 遍历当前场景的所有敌人
	if p.scene != nil {
		for _, enemy := range p.scene.GetArpgEnemies() {
			if enemy.IsDead() {
				continue
			}
			dist := p.distanceToEnemy(enemy)
			if dist < searchRange && dist < minDist {
				minDist = dist
				nearestEnemy = enemy
			}
		}
	}

	if nearestEnemy != nil {
		p.Arpg.TargetEnemy = nearestEnemy
	}
}

// performAttack 执行攻击
func (p *Role) performAttack() {
	p.SetAction(proto.RoleAction_RoleAction_AttackAxe)

	// 面向目标
	dx := p.Arpg.TargetEnemy.WX - p.GetWX()
	dy := p.Arpg.TargetEnemy.WY - (p.GetWY() - float32(p.sprite.roleImageSprite.Frame.H/2)) // 修正Y轴中心

	// 使用通用函数计算 8 方向朝向
	// 注意: 这里直接修改 p.sprite.orientation，与移动方向无关
	// 攻击时的朝向完全取决于目标相对于角色的位置
	p.sprite.orientation = commonct.CalculateOrientation(dx, dy)

	p.Arpg.TargetEnemy.TakeDamage(10)
}

// distanceToEnemy 计算到敌人的距离
func (p *Role) distanceToEnemy(enemy *ArpgEnemy) float32 {
	dx := p.GetWX() - enemy.WX
	dy := (p.GetWY() - float32(p.sprite.roleImageSprite.Frame.H/2)) - (enemy.WY - float32(enemy.GetCfg().Res.Move.FrameInfo[proto.AssetOrientation_AssetOrientation_Down][0].Frame.H/2))
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
