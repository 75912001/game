package user

import (
	"fmt"
	"image/color"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"saClient/src/cfg"
	"saClient/src/common"
	commoncamera "saClient/src/common/camera"
	commonrenderable "saClient/src/common/renderable"
	resfont "saClient/src/res/font"
)

func (p *Role) DrawAll(screen *ebitenv2.Image) {
	p.scene._map.DrawGround(screen, p.camera)

	renderables := p.collectRenderables()
	commonrenderable.SortYAndDraw(screen, p.camera, renderables)
	p.scene._map.RecycleTileSortInfoSlice(renderables) // 回收对象池

	p.scene._map.DrawOverhead(screen, p.camera)
	p.scene._map.DrawCollision(screen, p.camera)

	// 绘制过渡效果 (最上层)
	p.scene.transition.Draw(screen)

	// 绘制调试信息
	p.drawDebugInfo(screen, p.camera)
}

// GetWY 返回角色动作锚点的 World Y 坐标（用于Y-Sorting）
func (p *Role) GetWY() float32 {
	return p.sprite.wy
}

// Draw 仅绘制角色本身(实现 Renderable 接口)
func (p *Role) Draw(screen *ebitenv2.Image, cam *commoncamera.Camera) {
	screenX := p.sprite.centerWX - float32(cam.ViewportWX) - float32(p.sprite.roleImageSprite.Frame.W/2)
	screenY := p.sprite.centerWY - float32(cam.ViewportWY) - float32(p.sprite.roleImageSprite.Frame.H/2)

	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(p.sprite.image, op)

	// 绘制血条
	p.drawHPBar(screen, screenX, screenY, float32(p.sprite.roleImageSprite.Frame.W))
}

// drawHPBar 绘制角色血条
func (p *Role) drawHPBar(screen *ebitenv2.Image, x, y, w float32) {
	maxHP := p.BattleStats.GetHpMax()
	if maxHP == 0 {
		return
	}
	currentHP := p.GetCurrentHp()
	hpRatio := float64(currentHP) / float64(maxHP)
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
	opBg.ColorScale.ScaleWithColor(common.Colors_Gray)
	screen.DrawImage(whiteSubImage, opBg)

	// 绘制前景 (绿色)
	if hpRatio > 0 {
		opFg := &ebitenv2.DrawImageOptions{}
		opFg.GeoM.Scale(barW*hpRatio, barH)
		opFg.GeoM.Translate(barX, barY)
		opFg.ColorScale.ScaleWithColor(common.Colors_Green)
		screen.DrawImage(whiteSubImage, opFg)
	}

	// 绘制文字
	text := fmt.Sprintf("[hp:%d/%d]", currentHP, maxHP)
	textOp := &textv2.DrawOptions{}
	textOp.GeoM.Translate(barX, barY-14) // 文字位于血条上方
	textOp.ColorScale.ScaleWithColor(color.White)
	// 居中显示
	wText, _ := textv2.Measure(text, *resfont.GFace16, 0)
	textOp.GeoM.Translate((barW-wText)/2, 0)

	textv2.Draw(screen, text, *resfont.GFace16, textOp)
}

// collectRenderables 收集所有需要Y-Sorting的对象
func (p *Role) collectRenderables() []commonrenderable.IRenderable {
	// Y-Sorting 层
	renderables := make([]commonrenderable.IRenderable, 0, 1024)
	renderables = append(renderables, p) // 角色自己

	// 添加怪物
	for _, enemy := range p.scene._map.spawnManager.GetAllEnemies() {
		renderables = append(renderables, enemy)
	}

	renderables = p.scene._map.GetTileSortInfoSlice(p.camera, cfg.TiledLayerType_Building, renderables)
	renderables = p.scene._map.GetTileSortInfoSlice(p.camera, cfg.TiledLayerType_Objects, renderables)

	return renderables
}
