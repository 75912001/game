package user

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"saClient/src/cfg"
	"saClient/src/common"
	commoncamera "saClient/src/common/camera"
	"saClient/src/proto"
	resfont "saClient/src/res/font"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage    = ebitenv2.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebitenv2.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

// ArpgEnemy 怪物实体
type ArpgEnemy struct {
	SpawnPoint  *ArpgEnemySpawnPoint // 归属刷怪点
	Generated   *cfg.GeneratedEnemy  // 生成的敌人配置
	Level       uint32
	BattleStats *ArpgEnemyBattleStats // 战斗属性

	// 位置信息
	WX, WY      float32                // 世界坐标 (脚底中心)
	Orientation proto.AssetOrientation // 朝向

	AnimationFrame common.AnimationFrame

	AI *ArpgEnemyAI // AI 控制器
	HP uint32       // 当前生命值
}

// NewArpgEnemy 创建怪物实体
func NewArpgEnemy(spawnPoint *ArpgEnemySpawnPoint, generated *cfg.GeneratedEnemy, wx, wy float32) *ArpgEnemy {
	enemy := &ArpgEnemy{
		SpawnPoint:  spawnPoint,
		Generated:   generated,
		Level:       generated.Level,
		WX:          wx,
		WY:          wy,
		Orientation: proto.AssetOrientation_AssetOrientation_Down,
	}
	enemy.BattleStats = NewArpgEnemyBattleStats(enemy)
	enemy.AI = NewArpgEnemyAI(enemy)
	enemy.HP = enemy.BattleStats.GetHpMax()
	return enemy
}

func (p *ArpgEnemy) GetCfg() *cfg.Pet {
	return p.Generated.Pet
}

func (p *ArpgEnemy) IsDead() bool {
	return p.HP <= 0
}

// Update 每帧更新
func (p *ArpgEnemy) Update() {
	if p.IsDead() {
		return
	}

	p.AI.Update()
}

// GetWY 实现 IRenderable 接口 - 获取用于 Y-Sorting 的坐标
func (p *ArpgEnemy) GetWY() float32 {
	return p.WY
}

// Draw 实现 IRenderable 接口 - 绘制怪物
func (p *ArpgEnemy) Draw(screen *ebitenv2.Image, camera *commoncamera.Camera) {
	if p.IsDead() { // 死亡不绘制 todo menglc , 绘制死亡动画
		return
	}

	// 获取当前方向的动画帧
	frames := p.GetCfg().Res.Move.Frames[p.Orientation]
	img := frames[p.AnimationFrame.FrameIdx%uint32(len(frames))]
	// 获取帧信息
	frameInfos := p.GetCfg().Res.Move.FrameInfo[p.Orientation]
	frameInfo := frameInfos[p.AnimationFrame.FrameIdx%uint32(len(frameInfos))]
	// 计算屏幕坐标 (脚底中心为锚点)
	screenX := p.WX - float32(camera.ViewportWX)
	screenY := p.WY - float32(camera.ViewportWY)
	// 调整到图像左上角 (脚底中心 -> 图像左上角)
	screenX -= float32(frameInfo.Frame.W / 2)
	screenY -= float32(frameInfo.Frame.H)
	// 绘制
	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(img, op)

	// 绘制血条
	p.drawHPBar(screen, screenX, screenY, float32(frameInfo.Frame.W))
}

// drawHPBar 绘制血条
func (p *ArpgEnemy) drawHPBar(screen *ebitenv2.Image, x, y, w float32) {
	maxHP := p.BattleStats.GetHpMax()
	if maxHP == 0 {
		return
	}
	hpRatio := float64(p.HP) / float64(maxHP)
	if hpRatio < 0 {
		hpRatio = 0
	} else if hpRatio > 1 {
		hpRatio = 1
	}

	barW := float64(w)
	barH := 5.0
	barX := float64(x)
	barY := float64(y) - barH - 2 // 位于头顶上方 2 像素

	// 绘制背景 (灰色)
	opBg := &ebitenv2.DrawImageOptions{}
	opBg.GeoM.Scale(barW, barH)
	opBg.GeoM.Translate(barX, barY)
	opBg.ColorScale.ScaleWithColor(color.RGBA{100, 100, 100, 255})
	screen.DrawImage(whiteSubImage, opBg)

	// 绘制前景 (红色)
	if hpRatio > 0 {
		opFg := &ebitenv2.DrawImageOptions{}
		opFg.GeoM.Scale(barW*hpRatio, barH)
		opFg.GeoM.Translate(barX, barY)
		opFg.ColorScale.ScaleWithColor(color.RGBA{255, 0, 0, 255})
		screen.DrawImage(whiteSubImage, opFg)
	}

	// 绘制文字
	text := fmt.Sprintf("[lv:%d] [hp:%d/%d]", p.Level, p.HP, maxHP)
	textOp := &textv2.DrawOptions{}
	textOp.GeoM.Translate(barX, barY-14) // 文字位于血条上方
	textOp.ColorScale.ScaleWithColor(color.White)
	// 居中显示
	wText, _ := textv2.Measure(text, *resfont.GFace16, 0)
	textOp.GeoM.Translate((barW-wText)/2, 0)

	textv2.Draw(screen, text, *resfont.GFace16, textOp)
}

// DrawDebug 绘制调试信息 (视野范围)
func (p *ArpgEnemy) DrawDebug(screen *ebitenv2.Image, camera *commoncamera.Camera) {
	if p.IsDead() {
		return
	}

	// 计算屏幕坐标 (怪物中心)
	screenX := p.WX - float32(camera.ViewportWX)
	screenY := p.WY - float32(camera.ViewportWY)

	// 绘制视野范围 (红色虚线圆)
	colorView := color.RGBA{255, 0, 0, 200} // 红色
	drawDashedCircleEnemy(screen, screenX, screenY, cfg.GCommon.PetDefArpgViewRange, colorView, 48, 0.5, 2.0)
}

// drawDashedCircleEnemy 绘制虚线圆 (敌人专用)
func drawDashedCircleEnemy(screen *ebitenv2.Image, cx, cy, radius float32, clr color.RGBA, dashCount int, dashRatio float32, strokeWidth float32) {
	if dashCount <= 0 || radius <= 0 {
		return
	}

	angleStep := 2 * math.Pi / float64(dashCount)
	dashAngle := angleStep * float64(dashRatio)

	for i := 0; i < dashCount; i++ {
		startAngle := float64(i) * angleStep
		endAngle := startAngle + dashAngle

		// 计算起点和终点
		x1 := cx + radius*float32(math.Cos(startAngle))
		y1 := cy + radius*float32(math.Sin(startAngle))
		x2 := cx + radius*float32(math.Cos(endAngle))
		y2 := cy + radius*float32(math.Sin(endAngle))

		vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, clr, false)
	}
}

// TakeDamage 受到伤害
func (p *ArpgEnemy) TakeDamage(damage uint32) {
	if p.IsDead() {
		return
	}
	if damage <= p.HP {
		p.HP -= damage
	} else {
		p.HP = 0
	}
}
