package user

import ebitenv2 "github.com/hajimehoshi/ebiten/v2"

// Decoration 装饰物
type Decoration struct {
}

type DecorationMgr struct {
}

func NewDecorationMgr() *DecorationMgr {
	return &DecorationMgr{}
}

func (p *DecorationMgr) Update() {

}

func (p *DecorationMgr) Draw(screen *ebitenv2.Image) {

}
