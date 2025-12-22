package arpg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/cfg"
	commoncamera "saClient/src/common/camera"
	"saClient/src/proto"
)

// Enemy 怪物实体
type Enemy struct {
	SpawnPoint *EnemyGroupSpawnPoint // 归属刷怪点
	Generated  *cfg.GeneratedEnemy   // 生成的敌人配置

	// 位置信息
	WX, WY      float32 // 世界坐标 (脚底中心)
	TX, TY      float32 // Tile 坐标
	Orientation uint32  // 朝向

	// 动画信息
	FrameIdx  uint32 // 当前帧索引
	FrameTick uint32 // 帧计时器

	// AI
	AI *EnemyAI // AI 控制器

	// 战斗属性
	HP    int // 当前生命值
	MaxHP int // 最大生命值
	Level int // 等级

	// 状态
	IsDead bool // 是否死亡
}

// NewEnemy 创建怪物实体
func NewEnemy(spawnPoint *EnemyGroupSpawnPoint, generated *cfg.GeneratedEnemy, wx, wy float32) *Enemy {
	enemy := &Enemy{
		SpawnPoint:  spawnPoint,
		Generated:   generated,
		WX:          wx,
		WY:          wy,
		Orientation: uint32(proto.AssetOrientation_AssetOrientation_Down),
		Level:       int(generated.Level),
	}

	// 计算 Tile 坐标
	if spawnPoint != nil && spawnPoint.TiledMapCfg != nil {
		enemy.TX, enemy.TY = spawnPoint.TiledMapCfg.IsometricCT.W2T(wx, wy)
	}

	// 初始化 AI
	enemy.AI = NewEnemyAI(enemy)

	// 初始化属性 (简单计算)
	enemy.MaxHP = 100 + enemy.Level*10
	enemy.HP = enemy.MaxHP

	return enemy
}

func (p *Enemy) GetCfg() *cfg.Pet {
	return p.Generated.Config.Pet
}

// Update 每帧更新
func (p *Enemy) Update() {
	if p.IsDead {
		return
	}

	// 更新 AI
	if p.AI != nil {
		p.AI.Update()
	}
}

// GetWY 实现 IRenderable 接口 - 获取用于 Y-Sorting 的坐标
func (p *Enemy) GetWY() float32 {
	return p.WY
}

// Draw 实现 IRenderable 接口 - 绘制怪物
func (p *Enemy) Draw(screen *ebitenv2.Image, camera *commoncamera.Camera) {
	if p.IsDead {
		return
	}

	// 获取当前方向的动画帧
	frames := p.PetCfg.Res.Move.Frames[p.Orientation]
	if len(frames) == 0 {
		return
	}
	image := frames[p.FrameIdx%uint32(len(frames))]

	// 获取帧信息
	frameInfos := p.PetCfg.Res.Move.FrameInfo[p.Orientation]
	if len(frameInfos) == 0 {
		return
	}
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
func (p *Enemy) TakeDamage(damage int) {
	if p.IsDead {
		return
	}
	p.HP -= damage
	if p.HP <= 0 {
		p.HP = 0
		p.IsDead = true
	}
}

// UpdateAnimation 更新动画帧
func (p *Enemy) UpdateAnimation() {
	p.FrameTick++
	if p.FrameTick >= 6 { // 每 6 tick 切换一帧
		p.FrameTick = 0
		p.FrameIdx++
	}
}
