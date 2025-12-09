package game

import (
	"saClient/src/ui"
)

func (p *Game) Update() error {
	ui.GUIMgr.Update()
	p.user.Update()
	return nil
}
