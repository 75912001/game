package game

import (
	"saClient/src/ui"
)

func (p *Game) Update() error {
	ui.GUIMgr.Update()
	err := p.user.Update()
	if err != nil {
		return err
	}
	return nil
}
