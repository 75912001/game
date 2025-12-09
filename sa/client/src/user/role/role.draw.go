package role

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/common"
)

func (p *Role) Draw(screen *ebitenv2.Image) {
	// 绘制场景
	p.scene.Draw(screen, p.camera.X, p.camera.Y)

	// 绘制角色
	op := &ebitenv2.DrawImageOptions{}
	op.GeoM.Translate(float64(common.ScreenWidth/2-p.roleSprite.roleImageSprite.Frame.Width/2), float64(common.ScreenHeight/2-p.roleSprite.roleImageSprite.Frame.Height/2))
	screen.DrawImage(p.roleSprite.image, op)
}
