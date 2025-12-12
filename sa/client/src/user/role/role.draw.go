package role

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

func (p *Role) Draw(screen *ebitenv2.Image) {
	// 绘制场景
	p.scene.Draw(screen, p.camera)

	// 绘制角色
	// 角色屏幕位置 = 角色世界坐标 - 摄像机屏幕坐标 - 角色图片偏移
	screenX := p.roleSprite.centerTX - p.camera.ScreenX - p.roleSprite.roleImageSprite.Frame.Width/2
	screenY := p.roleSprite.centerTY - p.camera.ScreenY - p.roleSprite.roleImageSprite.Frame.Height/2

	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(p.roleSprite.image, op)
}
