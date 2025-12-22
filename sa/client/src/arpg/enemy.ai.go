package arpg

import (
	"math"
	"math/rand"
	"saClient/src/proto"
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
	Enemy        *Enemy       // 所属怪物
	State        EnemyAIState // 当前状态
	TargetWX     float32      // 目标点 X
	TargetWY     float32      // 目标点 Y
	StateTicks   int          // 当前状态已持续的 tick 数
	IdleDuration int          // Idle 状态持续时间
}

// NewEnemyAI 创建怪物 AI
func NewEnemyAI(enemy *Enemy) *EnemyAI {
	ai := &EnemyAI{
		Enemy:        enemy,
		State:        EnemyAIState_Idle,
		IdleDuration: 60 + rand.Intn(120), // 1-3 秒 (60 FPS)
	}
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
	p.StateTicks++
}

// updateIdle 待机状态更新
func (p *EnemyAI) updateIdle() {
	// Idle 超时后切换到 Patrol
	if p.StateTicks >= p.IdleDuration {
		p.switchToPatrol()
	}
}

// switchToPatrol 切换到巡逻状态
func (p *EnemyAI) switchToPatrol() {
	spawnPoint := p.Enemy.SpawnPoint
	if spawnPoint == nil || spawnPoint.TiledMapCfg == nil {
		return
	}

	mapCfg := spawnPoint.TiledMapCfg

	// 在巡逻半径内随机选择目标点
	angle := rand.Float64() * 2 * math.Pi
	distance := rand.Float32() * spawnPoint.Object.PatrolRadius

	targetWX := spawnPoint.WX + float32(math.Cos(angle))*distance
	targetWY := spawnPoint.WY + float32(math.Sin(angle))*distance

	// 限制在地图边界内
	targetWX, targetWY = mapCfg.ClampMapBoundaryWithW(targetWX, targetWY)

	// 检查目标点是否可达
	tx, ty := mapCfg.IsometricCT.W2T(targetWX, targetWY)
	if mapCfg.IsBlockedByTileWithTF(tx, ty) {
		// 目标不可达，重置 Idle 时间
		p.IdleDuration = 30 + rand.Intn(60) // 0.5-1.5 秒后重试
		p.StateTicks = 0
		return
	}

	// 设置目标并切换状态
	p.TargetWX = targetWX
	p.TargetWY = targetWY
	p.State = EnemyAIState_Patrol
	p.StateTicks = 0
}

// updatePatrol 巡逻状态更新
func (p *EnemyAI) updatePatrol() {
	enemy := p.Enemy
	spawnPoint := enemy.SpawnPoint
	if spawnPoint == nil || spawnPoint.TiledMapCfg == nil {
		p.switchToIdle()
		return
	}

	mapCfg := spawnPoint.TiledMapCfg

	// 计算移动方向
	dx := p.TargetWX - enemy.WX
	dy := p.TargetWY - enemy.WY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 检查是否到达目标
	if distance < 2.0 { // 阈值
		p.switchToIdle()
		return
	}

	// 归一化并移动
	moveSpeed := float32(1.5) // 怪物移动速度
	dx = dx / distance * moveSpeed
	dy = dy / distance * moveSpeed

	newWX := enemy.WX + dx
	newWY := enemy.WY + dy

	// 碰撞检测
	tx, ty := mapCfg.IsometricCT.W2T(newWX, newWY)
	if mapCfg.IsBlockedByTileWithTF(tx, ty) {
		// 碰到障碍物，回到 Idle
		p.switchToIdle()
		return
	}

	// 检查是否超出巡逻范围
	distToSpawn := float32(math.Sqrt(float64(
		(newWX-spawnPoint.WX)*(newWX-spawnPoint.WX) +
			(newWY-spawnPoint.WY)*(newWY-spawnPoint.WY))))
	if distToSpawn > spawnPoint.Object.PatrolRadius*1.5 { // 允许一定超出
		p.switchToIdle()
		return
	}

	// 更新位置
	enemy.WX = newWX
	enemy.WY = newWY
	enemy.TX = tx
	enemy.TY = ty

	// 更新朝向
	enemy.Orientation = p.calculateOrientation(dx, dy)

	// 更新动画帧
	enemy.UpdateAnimation()
}

// switchToIdle 切换到待机状态
func (p *EnemyAI) switchToIdle() {
	p.State = EnemyAIState_Idle
	p.StateTicks = 0
	p.IdleDuration = 60 + rand.Intn(120) // 1-3 秒
}

// calculateOrientation 根据移动方向计算朝向
func (p *EnemyAI) calculateOrientation(dx, dy float32) uint32 {
	// 8方向判断
	angle := math.Atan2(float64(dy), float64(dx))
	// 转换为 0-360 度
	degrees := angle * 180 / math.Pi
	if degrees < 0 {
		degrees += 360
	}

	// 根据角度确定方向
	// 右: -22.5 ~ 22.5, 右下: 22.5 ~ 67.5, 下: 67.5 ~ 112.5, 左下: 112.5 ~ 157.5
	// 左: 157.5 ~ 202.5, 左上: 202.5 ~ 247.5, 上: 247.5 ~ 292.5, 右上: 292.5 ~ 337.5
	switch {
	case degrees < 22.5 || degrees >= 337.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Right)
	case degrees < 67.5:
		return uint32(proto.AssetOrientation_AssetOrientation_DownRight)
	case degrees < 112.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Down)
	case degrees < 157.5:
		return uint32(proto.AssetOrientation_AssetOrientation_DownLeft)
	case degrees < 202.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Left)
	case degrees < 247.5:
		return uint32(proto.AssetOrientation_AssetOrientation_UpLeft)
	case degrees < 292.5:
		return uint32(proto.AssetOrientation_AssetOrientation_Up)
	default:
		return uint32(proto.AssetOrientation_AssetOrientation_UpRight)
	}
}
