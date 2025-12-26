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

// 追击相关常量
const (
	chasePathRecalcIntervalSec = 1    // 路径重计算间隔 (秒)
	chaseTargetMoveThreshold   = 50.0 // 目标移动阈值 (超过则重算路径)
	chaseAttackRange           = 30.0 // 攻击范围
	chaseArriveThreshold       = 5.0  // 到达路径点阈值
)

// ArpgEnemyAI 怪物 AI 控制器
type ArpgEnemyAI struct {
	Enemy       *ArpgEnemy       // 所属怪物
	State       ArpgEnemyAIState // 当前状态
	TargetWX    float32          // 目标点 X
	TargetWY    float32          // 目标点 Y
	IdleEndTime int64            // Idle 状态结束时间 (秒)

	// 追击寻路相关
	Pathfinder   *AStarPathfinder  // A* 寻路器
	Path         []AStarWorldPoint // 当前路径
	PathIndex    int               // 当前路径点索引
	LastPathTime int64             // 上次计算路径时间 (毫秒)
	LastTargetWX float32           // 上次目标位置 X
	LastTargetWY float32           // 上次目标位置 Y
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
	dx = dx / distance * cfg.GCommon.PetDefArpgMoveSpeed
	dy = dy / distance * cfg.GCommon.PetDefArpgMoveSpeed

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
	p.Path = nil
	p.PathIndex = 0
	p.LastPathTime = 0

	// 初始化寻路器 (如果地图支持逻辑网格)
	mapCfg := p.Enemy.SpawnPoint.TiledMapCfg
	if mapCfg.LogicalGrid != nil && p.Pathfinder == nil {
		p.Pathfinder = NewAStarPathfinder(mapCfg.LogicalGrid)
	}
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

	// 计算与目标的距离
	dx := targetWX - enemy.WX
	dy := targetWY - enemy.WY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 检查是否到达目标 (攻击范围)
	if distance < chaseAttackRange {
		// todo 切换到攻击状态
		return
	}

	// 使用 A* 寻路
	if p.Pathfinder != nil {
		p.updateChaseWithAStar(targetWX, targetWY, mapCfg)
	}
}

// updateChaseWithAStar 使用 A* 寻路追击
func (p *ArpgEnemyAI) updateChaseWithAStar(targetWX, targetWY float32, mapCfg *cfg.TiledMap) {
	enemy := p.Enemy
	now := xtime.GTimeMgr.ShadowTimestamp()

	// 检查是否需要重新计算路径
	needRecalc := false
	if p.Path == nil || len(p.Path) == 0 {
		needRecalc = true // 无路径
	} else if p.PathIndex >= len(p.Path) {
		needRecalc = true // 路径走完
	} else if now-p.LastPathTime >= chasePathRecalcIntervalSec {
		// 定时刷新: 检查目标是否移动超过阈值
		dx := targetWX - p.LastTargetWX
		dy := targetWY - p.LastTargetWY
		if dx*dx+dy*dy > chaseTargetMoveThreshold*chaseTargetMoveThreshold {
			needRecalc = true
		}
	}

	// 重新计算路径
	if needRecalc {
		p.Path = p.Pathfinder.FindPath(enemy.WX, enemy.WY, targetWX, targetWY)
		p.PathIndex = 0
		p.LastPathTime = now
		p.LastTargetWX = targetWX
		p.LastTargetWY = targetWY

		// 找不到路径, 等待下次重算
		if p.Path == nil || len(p.Path) == 0 {
			return
		}
	}

	// 沿路径移动
	if p.PathIndex < len(p.Path) {
		waypoint := p.Path[p.PathIndex]
		dx := waypoint.WX - enemy.WX
		dy := waypoint.WY - enemy.WY
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// 检查是否到达当前路径点
		if distance < chaseArriveThreshold {
			p.PathIndex++
			if p.PathIndex >= len(p.Path) {
				return // 路径走完
			}
			// 继续移动到下一个点
			waypoint = p.Path[p.PathIndex]
			dx = waypoint.WX - enemy.WX
			dy = waypoint.WY - enemy.WY
			distance = float32(math.Sqrt(float64(dx*dx + dy*dy)))
		}

		// 移动
		if distance > 0.01 {
			dx = dx / distance * cfg.GCommon.PetDefArpgMoveSpeed
			dy = dy / distance * cfg.GCommon.PetDefArpgMoveSpeed

			newWX := enemy.WX + dx
			newWY := enemy.WY + dy

			newWX, newWY, blocked := mapCfg.IsBlocked(newWX, newWY)
			if !blocked {
				enemy.WX = newWX
				enemy.WY = newWY
				enemy.Orientation = commonct.CalculateOrientation(dx, dy)
				enemy.UpdateAnimation()
			} else {
				// 路径上遇到阻挡 (动态障碍物?), 重算路径
				p.Path = nil
			}
		}
	}
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
