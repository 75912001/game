package user

import ebitenv2 "github.com/hajimehoshi/ebiten/v2"

// Plant 植物
type Plant struct {
}

type PlantMgr struct {
}

func NewPlantMgr() *PlantMgr {
	return &PlantMgr{}
}

func (p *PlantMgr) Update() {

}

func (p *PlantMgr) Draw(screen *ebitenv2.Image) {

}
