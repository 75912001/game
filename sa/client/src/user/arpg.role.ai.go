package user

import (
	"math"
	"saClient/src/cfg"
	commonct "saClient/src/common/coordinatetransform"
	"saClient/src/proto"

	xtime "github.com/75912001/xlib/time"
)

// ArpgRoleAIState 角色 arpg AI状态
type ArpgRoleAIState int

const (
	ArpgRoleAIState_Idle   ArpgRoleAIState = iota // 待机 (无自动行为)
	ArpgRoleAIState_Chase                         // 追击 (移动到攻击范围)
	ArpgRoleAIState_Attack                        // 攻击
	ArpgRoleAIState_Return                        // 归位 (返回记录位置)
)

// 角色战斗AI常量
const (
	arpgRoleAIPathRecalcIntervalSec int64   = 1    // 路径重计算间隔 (秒)
	arpgRoleAITargetMoveThreshold   float32 = 50.0 // 目标移动阈值 (超过则重算路径)
	arpgRoleAIArriveThreshold       float32 = 5.0  // 到达路径点阈值
	arpgRoleAIReturnArriveThreshold float32 = 1.0  // 归位到达阈值
)

// ArpgRoleAI 角色战斗AI控制器
type ArpgRoleAI struct {
	Role    *Role
	State   ArpgRoleAIState
	Enabled bool // 开关: 是否启用自动战斗

	// 归位点
	ReturnWX       float32
	ReturnWY       float32
	HasReturnPoint bool

	Target       *ArpgEnemy // 目标
	NextAttackMs int64      // 下一次攻击时间戳 (毫秒)

	// 攻击伤害延迟相关
	IsAttacking  bool       // 是否正在攻击动画中
	DamageDealt  bool       // 本次攻击是否已造成伤害
	AttackTarget *ArpgEnemy // 攻击锁定目标 (攻击开始时锁定)

	// 寻路相关
	Pathfinder   *AStarPathfinder  // A*寻路器
	Path         []AStarWorldPoint // 当前路径
	PathIndex    int               // 当前路径点索引
	LastPathTime int64             // 上次计算路径时间 (毫秒)
	LastTargetWX float32           // 上次目标位置 X
	LastTargetWY float32           // 上次目标位置 Y
}

// NewArpgRoleAI 创建角色战斗AI
func NewArpgRoleAI(role *Role) *ArpgRoleAI {
	return &ArpgRoleAI{
		Role:    role,
		State:   ArpgRoleAIState_Idle,
		Enabled: true, // 默认开启自动战斗
	}
}

// Update 每帧更新AI
func (p *ArpgRoleAI) Update() {
	if !p.Enabled {
		return
	}

	switch p.State {
	case ArpgRoleAIState_Idle:
		p.updateIdle()
	case ArpgRoleAIState_Chase:
		p.updateChase()
	case ArpgRoleAIState_Attack:
		p.updateAttack()
	case ArpgRoleAIState_Return:
		p.updateReturn()
	}
}

// updateIdle 待机状态: 检测视野内敌人
func (p *ArpgRoleAI) updateIdle() {
	target := p.selectTarget()
	if target != nil {
		// 记录当前位置为归位点
		p.recordReturnPoint()
		p.Target = target
		p.switchToChase()
	}
}

// updateChase 追击状态: A*寻路移动到攻击范围
func (p *ArpgRoleAI) updateChase() {
	// 检查目标是否有效
	if p.Target == nil || p.Target.IsDead() || !p.isTargetInVisionRange(p.Target) {
		// 目标无效或离开视野，尝试选择新目标
		newTarget := p.selectTarget()
		if newTarget != nil {
			p.Target = newTarget
		} else {
			p.switchToReturn()
			return
		}
	}

	// 检查是否在攻击范围内
	dist := p.distanceToTarget()
	attackRange := cfg.GCommon.GetRoleArpgDefAttackRangeByWeaponType(p.GetWeaponType())
	if dist <= attackRange {
		p.switchToAttack()
		return
	}

	// A*寻路移动
	p.moveToTarget()
}

// updateAttack 攻击状态: 执行攻击
func (p *ArpgRoleAI) updateAttack() {
	// 检查目标是否有效
	if p.Target == nil || p.Target.IsDead() || !p.isTargetInVisionRange(p.Target) {
		// 目标无效或离开视野，尝试选择新目标
		newTarget := p.selectTarget()
		if newTarget != nil {
			p.Target = newTarget
			// 检查新目标是否在攻击范围内
			dist := p.distanceToTarget()
			attackRange := cfg.GCommon.GetRoleArpgDefAttackRangeByWeaponType(p.GetWeaponType())
			if dist > attackRange {
				p.switchToChase()
				return
			}
		} else {
			p.switchToReturn()
			return
		}
	}

	// 检查目标是否逃出攻击范围
	dist := p.distanceToTarget()
	attackRange := cfg.GCommon.GetRoleArpgDefAttackRangeByWeaponType(p.GetWeaponType())
	if dist > attackRange {
		p.switchToChase()
		return
	}

	// 执行攻击 (调用现有攻击逻辑)
	p.performAttack()
}

// updateReturn 归位状态: 返回记录位置
func (p *ArpgRoleAI) updateReturn() {
	// 中途发现新敌人
	newTarget := p.selectTarget()
	if newTarget != nil {
		p.Target = newTarget
		p.switchToChase()
		return
	}

	// 检查是否到达归位点
	if !p.HasReturnPoint {
		p.switchToIdle()
		return
	}

	dx := p.ReturnWX - p.Role.GetWX()
	dy := p.ReturnWY - p.Role.GetWY()
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if dist <= arpgRoleAIReturnArriveThreshold {
		p.switchToIdle()
		return
	}

	// 移动到归位点
	p.moveToPoint(p.ReturnWX, p.ReturnWY)
}

// switchToIdle 切换到待机状态
func (p *ArpgRoleAI) switchToIdle() {
	p.State = ArpgRoleAIState_Idle
	p.Target = nil
	p.Path = nil
	p.PathIndex = 0
}

// switchToChase 切换到追击状态
func (p *ArpgRoleAI) switchToChase() {
	p.State = ArpgRoleAIState_Chase
	p.Path = nil
	p.PathIndex = 0
	p.LastPathTime = 0

	// 初始化寻路器
	p.initPathfinder()
}

// switchToAttack 切换到攻击状态
func (p *ArpgRoleAI) switchToAttack() {
	p.State = ArpgRoleAIState_Attack
	p.Path = nil
}

// switchToReturn 切换到归位状态
func (p *ArpgRoleAI) switchToReturn() {
	p.State = ArpgRoleAIState_Return
	p.Target = nil
	p.Path = nil
	p.PathIndex = 0
	p.LastPathTime = 0

	// 初始化寻路器
	p.initPathfinder()
}

// initPathfinder 初始化寻路器
func (p *ArpgRoleAI) initPathfinder() {
	if p.Role.scene == nil {
		return
	}
	mapCfg := p.Role.scene._map.tiledMapCfg
	if mapCfg.LogicalGrid != nil && p.Pathfinder == nil {
		p.Pathfinder = NewAStarPathfinder(mapCfg.LogicalGrid)
	}
}

// selectTarget 选择攻击目标 (视野内最近的敌人)
func (p *ArpgRoleAI) selectTarget() *ArpgEnemy {
	if p.Role.scene == nil {
		return nil
	}

	var nearestEnemy *ArpgEnemy
	var minDist float32 = math.MaxFloat32
	searchRange := cfg.GCommon.RoleArpgDefViewRange

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

	return nearestEnemy
}

// recordReturnPoint 记录归位点
func (p *ArpgRoleAI) recordReturnPoint() {
	p.ReturnWX = p.Role.GetWX()
	p.ReturnWY = p.Role.GetWY()
	p.HasReturnPoint = true
}

// UpdateReturnPoint 更新归位点 (手动移动时调用)
func (p *ArpgRoleAI) UpdateReturnPoint(wx, wy float32) {
	p.ReturnWX = wx
	p.ReturnWY = wy
	p.HasReturnPoint = true
}

// distanceToTarget 计算到当前目标的距离
func (p *ArpgRoleAI) distanceToTarget() float32 {
	if p.Target == nil {
		return math.MaxFloat32
	}
	return p.distanceToEnemy(p.Target)
}

// distanceToEnemy 计算到敌人的距离
func (p *ArpgRoleAI) distanceToEnemy(enemy *ArpgEnemy) float32 {
	roleWX := p.Role.GetWX()
	roleWY := p.Role.GetWY() - float32(p.Role.sprite.roleImageSprite.Frame.H/2) // 角色中心Y

	enemyWX := enemy.Sprite.GetWX()
	enemyWY := enemy.Sprite.GetWY() - float32(enemy.GetCfg().Res.Move.FrameInfo[proto.AssetOrientation_AssetOrientation_Down][0].Frame.H/2) // 敌人中心Y

	dx := roleWX - enemyWX
	dy := roleWY - enemyWY
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// isTargetInVisionRange 检查目标是否在视野范围内
func (p *ArpgRoleAI) isTargetInVisionRange(enemy *ArpgEnemy) bool {
	if enemy == nil {
		return false
	}
	dist := p.distanceToEnemy(enemy)
	return dist <= cfg.GCommon.RoleArpgDefViewRange
}

// moveToTarget 移动到目标 (A*寻路)
func (p *ArpgRoleAI) moveToTarget() {
	if p.Target == nil || p.Role.scene == nil {
		return
	}

	targetWX := p.Target.Sprite.GetWX()
	targetWY := p.Target.Sprite.GetWY()
	mapCfg := p.Role.scene._map.tiledMapCfg

	p.updatePathAndMove(targetWX, targetWY, mapCfg)
}

// moveToPoint 移动到指定点 (A*寻路)
func (p *ArpgRoleAI) moveToPoint(targetWX, targetWY float32) {
	if p.Role.scene == nil {
		return
	}
	mapCfg := p.Role.scene._map.tiledMapCfg
	p.updatePathAndMove(targetWX, targetWY, mapCfg)
}

// updatePathAndMove 更新路径并移动
func (p *ArpgRoleAI) updatePathAndMove(targetWX, targetWY float32, mapCfg *cfg.TiledMap) {
	now := xtime.GTimeMgr.ShadowTimestamp()

	// 检查是否需要重新计算路径
	needRecalc := false
	if p.Path == nil || len(p.Path) == 0 {
		needRecalc = true
	} else if p.PathIndex >= len(p.Path) {
		needRecalc = true
	} else if now-p.LastPathTime >= arpgRoleAIPathRecalcIntervalSec {
		dx := targetWX - p.LastTargetWX
		dy := targetWY - p.LastTargetWY
		if dx*dx+dy*dy > arpgRoleAITargetMoveThreshold*arpgRoleAITargetMoveThreshold {
			needRecalc = true
		}
	}

	// 重新计算路径
	if needRecalc && p.Pathfinder != nil {
		p.Path = p.Pathfinder.FindPath(p.Role.GetWX(), p.Role.GetWY(), targetWX, targetWY)
		p.PathIndex = 0
		p.LastPathTime = now
		p.LastTargetWX = targetWX
		p.LastTargetWY = targetWY

		if p.Path == nil || len(p.Path) == 0 {
			// 寻路失败，尝试直线移动
			p.moveDirectly(targetWX, targetWY, mapCfg)
			return
		}
	}

	// 沿路径移动
	if p.Path != nil && p.PathIndex < len(p.Path) {
		waypoint := p.Path[p.PathIndex]
		dx := waypoint.WX - p.Role.GetWX()
		dy := waypoint.WY - p.Role.GetWY()
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// 检查是否到达当前路径点
		if distance < arpgRoleAIArriveThreshold {
			p.PathIndex++
			if p.PathIndex >= len(p.Path) {
				return
			}
			waypoint = p.Path[p.PathIndex]
			dx = waypoint.WX - p.Role.GetWX()
			dy = waypoint.WY - p.Role.GetWY()
			distance = float32(math.Sqrt(float64(dx*dx + dy*dy)))
		}

		// 移动
		if distance > 0.01 {
			p.doMove(dx, dy, distance, mapCfg)
		}
	}
}

// moveDirectly 直线移动 (寻路失败时的备用方案)
func (p *ArpgRoleAI) moveDirectly(targetWX, targetWY float32, mapCfg *cfg.TiledMap) {
	dx := targetWX - p.Role.GetWX()
	dy := targetWY - p.Role.GetWY()
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance > 0.01 {
		p.doMove(dx, dy, distance, mapCfg)
	}
}

// doMove 执行移动
func (p *ArpgRoleAI) doMove(dx, dy, distance float32, mapCfg *cfg.TiledMap) {
	moveSpeed := cfg.GCommon.RoleDefMoveSpeed

	// 归一化方向
	dx = dx / distance * moveSpeed
	dy = dy / distance * moveSpeed

	newWX := p.Role.GetWX() + dx
	newWY := p.Role.GetWY() + dy

	// 边界限制
	newWX, newWY, blocked := mapCfg.Boundary.ClampWithW(newWX, newWY)
	if blocked {
		return
	}

	// Tile阻挡检测
	newTX, newTY := mapCfg.IsometricCT.W2T(newWX, newWY)
	if mapCfg.TileBlocked.IsBlockedWithTF(newTX, newTY) {
		return
	}

	// 对象阻挡检测
	if _, ok := mapCfg.FindBlockedByObject(newWX, newWY); ok {
		return
	}

	// 更新角色位置
	p.Role.sprite.wx = newWX
	p.Role.sprite.wy = newWY
	p.Role.sprite.tx, p.Role.sprite.ty = newTX, newTY

	// 更新朝向
	p.Role.sprite.orientation = commonct.CalculateOrientation(dx, dy)

	// 更新动画
	p.Role.SetAction(proto.RoleAction_RoleAction_Move)
	p.Role.animationFrame.Update()
	p.Role.UpdateWithAction()

	// 持久化位置
	p.Role.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WX, newWX)
	p.Role.SetValueF32(proto.AssetIDRecord_AssetIDRecord_BottomCenter_WY, newWY)
	p.Role.SetValueU64(proto.AssetIDRecord_AssetIDRecord_Orientation, uint64(p.Role.sprite.orientation))
}

// performAttack 执行攻击
func (p *ArpgRoleAI) performAttack() {
	// 检查攻击动画状态
	if p.IsAttacking {
		p.Role.animationFrame.Update()

		actionData := p.Role.cfgRole.ResRole.Actions[p.Role.sprite.action]
		frames := actionData.Frames[p.Role.sprite.orientation]
		frameInfos := actionData.FrameInfo[p.Role.sprite.orientation]
		currentFrameIdx := p.Role.animationFrame.FrameIdx % uint32(len(frames))

		// 检查当前帧是否为命中帧，且尚未造成伤害
		if !p.DamageDealt && currentFrameIdx < uint32(len(frameInfos)) {
			if frameInfos[currentFrameIdx].HitFrame {
				// 命中帧触发伤害
				if p.AttackTarget != nil && !p.AttackTarget.IsDead() {
					p.AttackTarget.TakeDamage(p.GetAttack())
				}
				p.DamageDealt = true
			}
		}

		// 检查动画是否播放完毕
		if p.Role.animationFrame.FrameIdx >= uint32(len(frames))-1 {
			p.Role.SetAction(proto.RoleAction_RoleAction_Move)
			p.IsAttacking = false
			p.DamageDealt = false
			p.AttackTarget = nil
			p.Role.UpdateWithAction()
		} else {
			p.Role.UpdateWithAction()
			return
		}
	}

	if p.Target == nil || p.Target.IsDead() {
		return
	}

	// 检查冷却
	nowMs := xtime.GTimeMgr.GetMillisecond()
	if p.NextAttackMs > nowMs {
		return
	}

	// 开始新的攻击
	p.Role.SetAction(proto.RoleAction_RoleAction_AttackAxe)

	// 面向目标
	dx := p.Target.Sprite.GetWX() - p.Role.GetWX()
	dy := p.Target.Sprite.GetWY() - (p.Role.GetWY() - float32(p.Role.sprite.roleImageSprite.Frame.H/2))
	p.Role.sprite.orientation = commonct.CalculateOrientation(dx, dy)

	// 设置攻击状态
	p.IsAttacking = true
	p.DamageDealt = false
	p.AttackTarget = p.Target // 锁定攻击目标

	// 设置冷却
	p.NextAttackMs = nowMs + cfg.GCommon.GetRoleArpgDefCdTimeMsByWeaponType(p.GetWeaponType())

	p.Role.UpdateWithAction()
}

// OnManualMove 手动移动时调用
func (p *ArpgRoleAI) OnManualMove(wx, wy float32) {
	// 更新归位点
	p.UpdateReturnPoint(wx, wy)

	// 如果在Return状态，中断归位
	if p.State == ArpgRoleAIState_Return {
		p.switchToIdle()
	}
}

// SetEnabled 设置是否启用自动战斗
func (p *ArpgRoleAI) SetEnabled(enabled bool) {
	p.Enabled = enabled
	if !enabled {
		p.switchToIdle()
	}
}

// IsAutoMoving 是否正在自动移动 (用于判断是否需要响应手动输入)
func (p *ArpgRoleAI) IsAutoMoving() bool {
	return p.Enabled && (p.State == ArpgRoleAIState_Chase || p.State == ArpgRoleAIState_Return)
}

// GetWeaponType 获取武器类型
func (p *ArpgRoleAI) GetWeaponType() proto.RoleWeaponType {
	return proto.RoleWeaponType_RoleWeaponType_Axe // todo menglc , 先写死斧头
}

// 获取攻击力
func (p *ArpgRoleAI) GetAttack() uint32 {
	return p.Role.BattleStats.GetAttack() + 10 // todo menglc , 先加个固定值,模拟装备攻击力
}
