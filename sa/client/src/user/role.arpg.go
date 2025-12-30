package user

import (
	xtime "github.com/75912001/xlib/time"
	"math"
	"saClient/src/cfg"
	commonct "saClient/src/common/coordinatetransform"
	"saClient/src/proto"
)

type RoleArpg struct {
	Role         *Role
	TargetEnemy  *ArpgEnemy // 攻击目标
	NextAttackMs int64      // 下一次攻击时间戳 (毫秒)
}

func NewRoleArpg(role *Role) *RoleArpg {
	return &RoleArpg{
		Role: role,
	}
}

// handleArpgAutoAttack 处理自动攻击逻辑
func (p *RoleArpg) handleAutoAttack() {
	// 1. 检查当前状态
	// 如果正在攻击且动画未播放完，保持攻击状态 (锁定)
	if p.Role.sprite.action == proto.RoleAction_RoleAction_AttackAxe {
		p.Role.animationFrame.Update()

		actionData := p.Role.cfgRole.ResRole.Actions[p.Role.sprite.action]
		frames := actionData.Frames[p.Role.sprite.orientation]
		// 简单判断：如果帧索引超过了总帧数，说明播放了一轮
		if p.Role.animationFrame.FrameIdx >= uint32(len(frames))-1 { // 检查动画播放完毕
			p.Role.SetAction(proto.RoleAction_RoleAction_Move) // 攻击结束，切回默认状态
			p.Role.UpdateWithAction()                          // 更新图像
		} else { // 继续播放攻击动画
			p.Role.UpdateWithAction() // 更新图像
			return
		}
	}

	p.selectTarget()
	if p.TargetEnemy == nil { // 没有目标
		return
	}
	// 检查距离
	dist := p.distanceToEnemy(p.TargetEnemy)

	attackRange := cfg.GCommon.GetRoleArpgDefAttackRangeByWeaponType(proto.RoleWeaponType_RoleWeaponType_Axe)
	if dist <= attackRange { // 目标在范围内
		// 检查冷却
		nowMillisecond := xtime.GTimeMgr.GetMillisecond()
		if p.NextAttackMs <= nowMillisecond { // 冷却结束
			p.performAttack()
			p.NextAttackMs = nowMillisecond + cfg.GCommon.GetRoleArpgDefCdTimeMsByWeaponType(proto.RoleWeaponType_RoleWeaponType_Axe)
			p.Role.UpdateWithAction() // 更新图像
		}
	} else { // 目标在范围外
		// todo menglc 添加自动追击逻辑，或者什么都不做保持原地
	}
}

// selectTarget 选择攻击目标
func (p *RoleArpg) selectTarget() {
	if p.TargetEnemy != nil { // 已有目标
		if p.TargetEnemy.IsDead() { // 目标死亡
			p.TargetEnemy = nil
		} else { // 目标存活
			return
		}
	}
	// 寻找最近的敌人
	var nearestEnemy *ArpgEnemy
	var minDist float32 = math.MaxFloat32 // 最小距离
	searchRange := cfg.GCommon.RoleArpgDefViewRange

	// 遍历当前场景的所有敌人
	if p.Role.scene != nil {
		for _, enemy := range p.Role.scene.GetArpgEnemies() {
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
		p.TargetEnemy = nearestEnemy
	}
}

// performAttack 执行攻击
func (p *RoleArpg) performAttack() {
	p.Role.SetAction(proto.RoleAction_RoleAction_AttackAxe)

	// 面向目标
	dx := p.TargetEnemy.WX - p.Role.GetWX()
	dy := p.TargetEnemy.WY - (p.Role.GetWY() - float32(p.Role.sprite.roleImageSprite.Frame.H/2)) // 修正Y轴中心

	// 使用通用函数计算 8 方向朝向
	p.Role.sprite.orientation = commonct.CalculateOrientation(dx, dy)

	p.TargetEnemy.TakeDamage(p.GetAttack())
}

// distanceToEnemy 计算到敌人的距离
func (p *RoleArpg) distanceToEnemy(enemy *ArpgEnemy) float32 {
	dx := p.Role.GetWX() - enemy.WX
	dy := (p.Role.GetWY() - float32(p.Role.sprite.roleImageSprite.Frame.H/2)) - (enemy.WY - float32(enemy.GetCfg().Res.Move.FrameInfo[proto.AssetOrientation_AssetOrientation_Down][0].Frame.H/2))
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// GetWeaponType 获取武器类型
func (p *RoleArpg) GetWeaponType() proto.RoleWeaponType {
	return proto.RoleWeaponType_RoleWeaponType_Axe // todo menglc , 先写死斧头
}

// 获取攻击力
func (p *RoleArpg) GetAttack() uint32 {
	return p.Role.BattleStats.GetAttack() + 10 // todo menglc , 先加个固定值,模拟装备攻击力
}
