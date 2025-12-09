package role

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/proto"
)

func (p *Role) Draw(screen *ebitenv2.Image) {
	direction := p.GetValueU32(proto.AssetIDRecord_AssetIDRecord_Direction)
	frames := p.cfgRole.ResRole.Move.Frames[direction]
	frameImage := frames[p.frameIdx%uint32(len(frames))]
	// 绘制
	op := &ebitenv2.DrawImageOptions{}
	x := p.GetValueU32(proto.AssetIDRecord_AssetIDRecord_X)
	y := p.GetValueU32(proto.AssetIDRecord_AssetIDRecord_Y)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(frameImage, op)
}
