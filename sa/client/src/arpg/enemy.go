package arpg

import (
	"saClient/src/cfg"
	"saClient/src/common"
	commoncamera "saClient/src/common/camera"
	"saClient/src/proto"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// Enemy 怪物实体
type Enemy struct {
	SpawnPoint  *EnemySpawnPoint    // 归属刷怪点
	Generated   *cfg.GeneratedEnemy // 生成的敌人配置
	Level       uint32
	BattleStats *EnemyBattleStats // 战斗属性

	// 位置信息
	WX, WY      float32 // 世界坐标 (脚底中心)
	Orientation uint32  // 朝向

	// 动画信息
	FrameIdx  uint32 // 当前帧索引
	FrameTick uint32 // 帧计时器

	AI *EnemyAI // AI 控制器
	HP uint32   // 当前生命值
}

// NewEnemy 创建怪物实体
func NewEnemy(spawnPoint *EnemySpawnPoint, generated *cfg.GeneratedEnemy, wx, wy float32) *Enemy {
	enemy := &Enemy{
		SpawnPoint:  spawnPoint,
		Generated:   generated,
		Level:       generated.Level,
		WX:          wx,
		WY:          wy,
		Orientation: uint32(proto.AssetOrientation_AssetOrientation_Down),
		FrameIdx:    0,
		FrameTick:   0,
	}
	enemy.BattleStats = NewEnemyBattleStats(enemy)
	enemy.AI = NewEnemyAI(enemy)
	enemy.HP = enemy.BattleStats.GetHpMax()
	return enemy
}

func (p *Enemy) GetCfg() *cfg.Pet {
	return p.Generated.Config.Pet
}

func (p *Enemy) IsDead() bool {
	return p.HP <= 0
}

// Update 每帧更新
func (p *Enemy) Update() {
	if p.IsDead() {
		return
	}

	p.AI.Update()
}

// GetWY 实现 IRenderable 接口 - 获取用于 Y-Sorting 的坐标
func (p *Enemy) GetWY() float32 {
	return p.WY
}

// Draw 实现 IRenderable 接口 - 绘制怪物
func (p *Enemy) Draw(screen *ebitenv2.Image, camera *commoncamera.Camera) {
	if p.IsDead() { // 死亡不绘制 todo menglc , 绘制死亡动画
		return
	}

	// 获取当前方向的动画帧
	frames := p.GetCfg().Res.Move.Frames[p.Orientation]
	image := frames[p.FrameIdx%uint32(len(frames))]
	// 获取帧信息
	frameInfos := p.GetCfg().Res.Move.FrameInfo[p.Orientation]
	frameInfo := frameInfos[p.FrameIdx%uint32(len(frameInfos))]
	// 计算屏幕坐标 (脚底中心为锚点)
	screenX := p.WX - float32(camera.ViewportWX)
	screenY := p.WY - float32(camera.ViewportWY)
	// 调整到图像左上角 (脚底中心 -> 图像左上角)
	screenX -= float32(frameInfo.Frame.Width / 2)
	screenY -= float32(frameInfo.Frame.Height)
	// 绘制
	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(image, op)
}

// TakeDamage 受到伤害
func (p *Enemy) TakeDamage(damage uint32) {
	if p.IsDead() {
		return
	}
	if damage <= p.HP {
		p.HP -= damage
	} else {
		p.HP = 0
	}
}

// UpdateAnimation 更新动画帧
func (p *Enemy) UpdateAnimation() {
	p.FrameTick++
	if common.FrameTickPerChange <= p.FrameTick {
		p.FrameTick = 0
		p.FrameIdx++
	}
}
