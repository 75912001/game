package user

import (
	xtime "github.com/75912001/xlib/time"
	xutil "github.com/75912001/xlib/util"
	"math"
	"saClient/src/cfg"
	commonct "saClient/src/common/coordinatetransform"
	"saClient/src/proto"
)

// ArpgEnemyAIState 怪物 AI 状态
type ArpgEnemyAIState int

const (
	ArpgEnemyAIState_Idle   ArpgEnemyAIState = iota // 待机
	ArpgEnemyAIState_Patrol                         // 巡逻
	ArpgEnemyAIState_Chase                          // 追击
	ArpgEnemyAIState_Attack                         // 攻击
)

// 追击相关常量
const (
	arpgEnemyChasePathRecalcIntervalSec int64   = 1    // 路径重计算间隔 (秒)
	arpgEnemyChaseTargetMoveThreshold   float32 = 50.0 // 目标移动阈值 (超过则重算路径)
	arpgEnemyChaseArriveThreshold       float32 = 5.0  // 到达路径点阈值
)

// ArpgEnemyAI 怪物 AI 控制器
type ArpgEnemyAI struct {
	Enemy       *ArpgEnemy       // 所属怪物
	State       ArpgEnemyAIState // 当前状态
	TargetWX    float32          // 目标点 X
	TargetWY    float32          // 目标点 Y
	IdleEndTime int64            // Idle 状态结束时间 (秒)

	// 攻击相关
	IsAttacking  bool  // 是否正在攻击动画中
	DamageDealt  bool  // 本次攻击是否已造成伤害
	NextAttackMs int64 // 下一次攻击时间戳 (毫秒)

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
	case ArpgEnemyAIState_Attack:
		p.updateAttack()
	}
}

// updateIdle 待机状态更新
func (p *ArpgEnemyAI) updateIdle() {
	if p.isTargetInChaseRange() { // 目标在追击范围内
		if p.isTargetInVisionRange() { // 目标在视野范围内
			p.switchToChase()
			return
		}
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
	if p.isTargetInChaseRange() { // 目标在追击范围内
		if p.isTargetInVisionRange() { // 目标在视野范围内
			p.switchToChase()
			return
		}
	}

	enemy := p.Enemy
	spawnPoint := enemy.SpawnPoint
	mapCfg := spawnPoint.TiledMapCfg

	// 计算移动方向
	dx := p.TargetWX - enemy.Sprite.GetWX()
	dy := p.TargetWY - enemy.Sprite.GetWY()
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 检查是否到达目标
	if xutil.Float32Less(distance, 1.0) { // 阈值
		p.switchToIdle()
		return
	}

	// 移动
	dx = dx / distance * cfg.GCommon.PetDefArpgMoveSpeed
	dy = dy / distance * cfg.GCommon.PetDefArpgMoveSpeed

	newWX := enemy.Sprite.GetWX() + dx
	newWY := enemy.Sprite.GetWY() + dy

	newWX, newWY, blocked := mapCfg.IsBlocked(newWX, newWY)
	if blocked { // 阻挡
		p.switchToIdle()
		return
	}

	// 更新位置和中心点
	enemy.SetPosition(newWX, newWY)

	// 更新朝向
	enemy.Sprite.SetOrientation(commonct.CalculateOrientation(dx, dy))

	// 更新动画帧
	enemy.AnimationFrame.Update()
}

// switchToIdle 切换到待机状态
func (p *ArpgEnemyAI) switchToIdle() {
	p.State = ArpgEnemyAIState_Idle
	p.IdleEndTime = xtime.GTimeMgr.ShadowTimestamp() + xutil.RandomInt64(1, 2)
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
	if !p.isTargetInChaseRange() { // 目标不在追击范围内
		p.switchToPatrol()
		return
	}

	enemy := p.Enemy
	spawnPoint := enemy.SpawnPoint
	mapCfg := spawnPoint.TiledMapCfg

	target := GUser.role
	targetWX := target.GetWX()
	targetWY := target.GetWY()
	// 计算与目标的距离
	dx := targetWX - enemy.Sprite.GetWX()
	dy := targetWY - enemy.Sprite.GetWY()
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 检查是否到达目标 (攻击范围)
	if distance < cfg.GCommon.PetArpgDefAttackRange {
		p.switchToAttack()
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
	} else if now-p.LastPathTime >= arpgEnemyChasePathRecalcIntervalSec {
		// 定时刷新: 检查目标是否移动超过阈值
		dx := targetWX - p.LastTargetWX
		dy := targetWY - p.LastTargetWY
		if dx*dx+dy*dy > arpgEnemyChaseTargetMoveThreshold*arpgEnemyChaseTargetMoveThreshold {
			needRecalc = true
		}
	}

	// 重新计算路径
	if needRecalc {
		p.Path = p.Pathfinder.FindPath(enemy.Sprite.GetWX(), enemy.Sprite.GetWY(), targetWX, targetWY)
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
		dx := waypoint.WX - enemy.Sprite.GetWX()
		dy := waypoint.WY - enemy.Sprite.GetWY()
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// 检查是否到达当前路径点
		if distance < arpgEnemyChaseArriveThreshold {
			p.PathIndex++
			if p.PathIndex >= len(p.Path) {
				return // 路径走完
			}
			// 继续移动到下一个点
			waypoint = p.Path[p.PathIndex]
			dx = waypoint.WX - enemy.Sprite.GetWX()
			dy = waypoint.WY - enemy.Sprite.GetWY()
			distance = float32(math.Sqrt(float64(dx*dx + dy*dy)))
		}

		// 移动
		if distance > 0.01 {
			dx = dx / distance * cfg.GCommon.PetDefArpgMoveSpeed
			dy = dy / distance * cfg.GCommon.PetDefArpgMoveSpeed

			newWX := enemy.Sprite.GetWX() + dx
			newWY := enemy.Sprite.GetWY() + dy

			newWX, newWY, blocked := mapCfg.IsBlocked(newWX, newWY)
			if !blocked {
				enemy.SetPosition(newWX, newWY)
				enemy.Sprite.SetOrientation(commonct.CalculateOrientation(dx, dy))
				enemy.AnimationFrame.Update()
			} else {
				// 路径上遇到阻挡 (动态障碍物?), 重算路径
				p.Path = nil
			}
		}
	}
}

// 检查目标是否在视野范围内
func (p *ArpgEnemyAI) isTargetInVisionRange() bool {
	target := GUser.role

	// 检查距离
	dx := target.GetWX() - p.Enemy.Sprite.GetWX()
	dy := target.GetWY() - p.Enemy.Sprite.GetWY()
	distSq := dx*dx + dy*dy
	return distSq <= cfg.GCommon.PetDefArpgViewRange*cfg.GCommon.PetDefArpgViewRange
}

// 检查目标是否在追击范围内
func (p *ArpgEnemyAI) isTargetInChaseRange() bool {
	target := GUser.role
	spawnPoint := p.Enemy.SpawnPoint

	distToSpawn := float32(math.Sqrt(float64(
		(target.GetWX()-spawnPoint.Object.WX)*(target.GetWX()-spawnPoint.Object.WX) +
			(target.GetWY()-spawnPoint.Object.WY)*(target.GetWY()-spawnPoint.Object.WY))))

	return distToSpawn <= spawnPoint.Object.PatrolRadius*cfg.GCommon.PetDefArpgChaseRadiusMultiplier
}

// switchToAttack 切换到攻击状态
func (p *ArpgEnemyAI) switchToAttack() {
	p.State = ArpgEnemyAIState_Attack
	p.IsAttacking = false
	p.DamageDealt = false
}

// updateAttack 攻击状态更新
func (p *ArpgEnemyAI) updateAttack() {
	enemy := p.Enemy
	target := GUser.role

	// 检查目标是否离开追击范围
	if !p.isTargetInChaseRange() {
		p.IsAttacking = false
		p.DamageDealt = false
		enemy.SetAction(proto.PetAction_PetAction_Move)
		p.switchToPatrol()
		return
	}

	// 计算与目标的距离
	dx := target.GetWX() - enemy.Sprite.GetWX()
	dy := target.GetWY() - enemy.Sprite.GetWY()
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 如果目标逃出攻击范围，切回追击
	if distance > cfg.GCommon.PetArpgDefAttackRange {
		p.IsAttacking = false
		p.DamageDealt = false
		enemy.SetAction(proto.PetAction_PetAction_Move)
		p.switchToChase()
		return
	}

	// 检查攻击动画状态
	if p.IsAttacking {
		enemy.AnimationFrame.Update()

		// 获取攻击动画帧信息
		attackData := enemy.GetCfg().Res.Attack
		if attackData == nil {
			// 没有攻击动画，直接造成伤害并切回待机
			GUser.role.TakeDamage(enemy.BattleStats.GetAttack())
			p.switchToIdle()
			return
		}

		orientation := enemy.Sprite.GetOrientation()
		frames := attackData.Frames[orientation]
		frameInfos := attackData.FrameInfo[orientation]
		if len(frames) == 0 {
			// 没有该方向的攻击动画
			p.IsAttacking = false
			p.switchToIdle()
			return
		}

		currentFrameIdx := enemy.AnimationFrame.FrameIdx % uint32(len(frames))

		// 检查当前帧是否为命中帧，且尚未造成伤害
		if !p.DamageDealt && currentFrameIdx < uint32(len(frameInfos)) {
			if frameInfos[currentFrameIdx].HitFrame {
				// 命中帧触发伤害
				GUser.role.TakeDamage(enemy.BattleStats.GetAttack())
				p.DamageDealt = true
			}
		}

		// 检查动画是否播放完毕
		if enemy.AnimationFrame.FrameIdx >= uint32(len(frames))-1 {
			enemy.SetAction(proto.PetAction_PetAction_Move)
			p.IsAttacking = false
			p.DamageDealt = false
		}
		return
	}

	// 检查冷却
	nowMs := xtime.GTimeMgr.GetMillisecond()
	if p.NextAttackMs > nowMs {
		return
	}

	// 开始新的攻击
	enemy.SetAction(proto.PetAction_PetAction_Attack)
	enemy.AnimationFrame.Reset()

	// 面向目标
	enemy.Sprite.SetOrientation(commonct.CalculateOrientation(dx, dy))

	// 设置攻击状态
	p.IsAttacking = true
	p.DamageDealt = false

	// 设置冷却
	p.NextAttackMs = nowMs + cfg.GCommon.PetArpgDefCdTimeMs
}
