package user

import (
	xtime "github.com/75912001/xlib/time"
	xutil "github.com/75912001/xlib/util"
	"math"
	"saClient/src/cfg"
	commonct "saClient/src/common/coordinatetransform"
)

// ArpgEnemyAIState 怪物 AI 状态
type ArpgEnemyAIState int

const (
	ArpgEnemyAIState_Idle   ArpgEnemyAIState = iota // 待机
	ArpgEnemyAIState_Patrol                         // 巡逻
	ArpgEnemyAIState_Chase                          // 追击
	// 后续扩展
	// ArpgEnemyAIState_Attack                  // 攻击
)

// ArpgEnemyAI 怪物 AI 控制器
type ArpgEnemyAI struct {
	Enemy       *ArpgEnemy       // 所属怪物
	State       ArpgEnemyAIState // 当前状态
	TargetWX    float32          // 目标点 X
	TargetWY    float32          // 目标点 Y
	IdleEndTime int64            // Idle 状态结束时间 (秒)
}

// NewArpgEnemyAI 创建怪物 AI
func NewArpgEnemyAI(enemy *ArpgEnemy) *ArpgEnemyAI {
	ai := &ArpgEnemyAI{
		Enemy: enemy,
	}
	ai.switchToPatrol()
	return ai
}

// Update 每帧更新 AI
func (p *ArpgEnemyAI) Update() {
	switch p.State {
	case ArpgEnemyAIState_Idle:
		p.updateIdle()
	case ArpgEnemyAIState_Patrol:
		p.updatePatrol()
	case ArpgEnemyAIState_Chase:
		p.updateChase()
	}
}

// updateIdle 待机状态更新
func (p *ArpgEnemyAI) updateIdle() {
	// 检查视野
	if p.checkVision() {
		return
	}

	// Idle 超时后切换到 Patrol
	if xtime.GTimeMgr.ShadowTimestamp() >= p.IdleEndTime {
		p.switchToPatrol()
	}
}

// switchToPatrol 切换到巡逻状态
func (p *ArpgEnemyAI) switchToPatrol() {
	spawnPoint := p.Enemy.SpawnPoint

	// 在巡逻半径内随机选择目标点
	targetWX, targetWY := spawnPoint.RandomPositionInRadius(spawnPoint.Object.PatrolRadius)

	// 设置目标并切换状态
	p.TargetWX = targetWX
	p.TargetWY = targetWY
	p.State = ArpgEnemyAIState_Patrol
}

// updatePatrol 巡逻状态更新
func (p *ArpgEnemyAI) updatePatrol() {
	// 检查视野
	if p.checkVision() {
		return
	}

	enemy := p.Enemy
	spawnPoint := enemy.SpawnPoint
	mapCfg := spawnPoint.TiledMapCfg

	// 计算移动方向
	dx := p.TargetWX - enemy.WX
	dy := p.TargetWY - enemy.WY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 检查是否到达目标
	if xutil.Float32Less(distance, 1.0) { // 阈值
		p.switchToIdle()
		return
	}

	// 移动
	dx = dx / distance * cfg.GCommon.PetDefArpgSpeed
	dy = dy / distance * cfg.GCommon.PetDefArpgSpeed

	newWX := enemy.WX + dx
	newWY := enemy.WY + dy

	newWX, newWY, blocked := mapCfg.IsBlocked(newWX, newWY)
	if blocked { // 阻挡
		p.switchToIdle()
		return
	}

	// 更新位置
	enemy.WX = newWX
	enemy.WY = newWY

	// 更新朝向
	enemy.Orientation = commonct.CalculateOrientation(dx, dy)

	// 更新动画帧
	enemy.UpdateAnimation()
}

// switchToIdle 切换到待机状态
func (p *ArpgEnemyAI) switchToIdle() {
	p.State = ArpgEnemyAIState_Idle
	p.IdleEndTime = xtime.GTimeMgr.ShadowTimestamp() + xutil.RandomInt64(1, 1)
}

// switchToChase 切换到追击状态
func (p *ArpgEnemyAI) switchToChase() {
	p.State = ArpgEnemyAIState_Chase
}

// updateChase 追击状态更新
func (p *ArpgEnemyAI) updateChase() {
	enemy := p.Enemy
	spawnPoint := enemy.SpawnPoint
	mapCfg := spawnPoint.TiledMapCfg

	target := GUser.role
	targetWX := target.GetWX()
	targetWY := target.GetWY()

	// 检查是否超出追击范围 (2倍巡逻半径)
	distToSpawn := float32(math.Sqrt(float64(
		(targetWX-spawnPoint.Object.WX)*(targetWX-spawnPoint.Object.WX) +
			(targetWY-spawnPoint.Object.WY)*(targetWY-spawnPoint.Object.WY))))

	if distToSpawn > spawnPoint.Object.PatrolRadius*2 {
		p.switchToPatrol()
		return
	}

	// 计算移动方向
	dx := targetWX - enemy.WX
	dy := targetWY - enemy.WY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 检查是否到达目标 (攻击范围)
	if xutil.Float32Less(distance, 30.0) { // 假设攻击范围 30
		// todo 切换到攻击状态
		return
	}

	// 移动
	dx = dx / distance * cfg.GCommon.PetDefArpgSpeed
	dy = dy / distance * cfg.GCommon.PetDefArpgSpeed

	newWX := enemy.WX + dx
	newWY := enemy.WY + dy

	newWX, newWY, blocked := mapCfg.IsBlocked(newWX, newWY)
	// 追击时遇到阻挡，简单处理为停止移动，或者尝试滑动（这里简化处理）
	if blocked {
		// 简单的滑动处理：尝试只移动 X 或只移动 Y
		newWX2, _, blockedX := mapCfg.IsBlocked(enemy.WX+dx, enemy.WY)
		if !blockedX {
			newWX = newWX2
			newWY = enemy.WY
		} else {
			_, newWY2, blockedY := mapCfg.IsBlocked(enemy.WX, enemy.WY+dy)
			if !blockedY {
				newWX = enemy.WX
				newWY = newWY2
			} else {
				// 完全阻挡，停止
				return
			}
		}
	}

	// 更新位置
	enemy.WX = newWX
	enemy.WY = newWY

	// 更新朝向
	enemy.Orientation = commonct.CalculateOrientation(dx, dy)

	// 更新动画帧
	enemy.UpdateAnimation()
}

// checkVision 检查视野
func (p *ArpgEnemyAI) checkVision() bool {
	target := GUser.role

	// 检查距离
	dx := target.GetWX() - p.Enemy.WX
	dy := target.GetWY() - p.Enemy.WY
	distSq := dx*dx + dy*dy
	viewRange := cfg.GCommon.PetDefViewRange

	if distSq <= viewRange*viewRange {
		p.switchToChase()
		return true
	}
	return false
}
