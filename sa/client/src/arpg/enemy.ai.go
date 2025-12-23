package arpg

import (
	xtime "github.com/75912001/xlib/time"
	xutil "github.com/75912001/xlib/util"
	"math"
	commonct "saClient/src/common/coordinatetransform"
)

// EnemyAIState 怪物 AI 状态
type EnemyAIState int

const (
	EnemyAIState_Idle   EnemyAIState = iota // 待机
	EnemyAIState_Patrol                     // 巡逻
	// 后续扩展
	// EnemyAIState_Chase                   // 追击
	// EnemyAIState_Attack                  // 攻击
)

// EnemyAI 怪物 AI 控制器
type EnemyAI struct {
	Enemy       *Enemy       // 所属怪物
	State       EnemyAIState // 当前状态
	TargetWX    float32      // 目标点 X
	TargetWY    float32      // 目标点 Y
	IdleEndTime int64        // Idle 状态结束时间 (秒)
}

// NewEnemyAI 创建怪物 AI
func NewEnemyAI(enemy *Enemy) *EnemyAI {
	ai := &EnemyAI{
		Enemy: enemy,
	}
	ai.switchToPatrol()
	return ai
}

// Update 每帧更新 AI
func (p *EnemyAI) Update() {
	switch p.State {
	case EnemyAIState_Idle:
		p.updateIdle()
	case EnemyAIState_Patrol:
		p.updatePatrol()
	}
}

// updateIdle 待机状态更新
func (p *EnemyAI) updateIdle() {
	// Idle 超时后切换到 Patrol
	if xtime.GTimeMgr.ShadowTimestamp() >= p.IdleEndTime {
		p.switchToPatrol()
	}
}

// switchToPatrol 切换到巡逻状态
func (p *EnemyAI) switchToPatrol() {
	spawnPoint := p.Enemy.SpawnPoint

	// 在巡逻半径内随机选择目标点
	targetWX, targetWY := spawnPoint.RandomPositionInRadius(spawnPoint.Object.PatrolRadius)

	// 设置目标并切换状态
	p.TargetWX = targetWX
	p.TargetWY = targetWY
	p.State = EnemyAIState_Patrol
}

// updatePatrol 巡逻状态更新
func (p *EnemyAI) updatePatrol() {
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
	moveSpeed := enemy.Generated.Config.Pet.Attributes.ArpgSpeed
	dx = dx / distance * moveSpeed
	dy = dy / distance * moveSpeed

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
	enemy.TX, enemy.TY = mapCfg.IsometricCT.W2T(newWX, newWY)

	// 更新朝向
	enemy.Orientation = commonct.CalculateOrientation(dx, dy)

	// 更新动画帧
	enemy.UpdateAnimation()
}

// switchToIdle 切换到待机状态
func (p *EnemyAI) switchToIdle() {
	p.State = EnemyAIState_Idle
	p.IdleEndTime = xtime.GTimeMgr.ShadowTimestamp() + xutil.RandomInt64(1, 1)
}
