package game

import (
	xtime "github.com/75912001/xlib/time"
	"saClient/src/ui"
)

func (p *Game) Update() error {
	xtime.GTimeMgr.Update()
	ui.GUIMgr.Update()
	p.user.Update()
	return nil
}
