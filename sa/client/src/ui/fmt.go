package ui

import (
	"fmt"
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	resfont "saClient/src/res/font"
)

func Printf(screen *ebitenv2.Image, x, y float64, format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	textOp := &textv2.DrawOptions{}
	textOp.GeoM.Translate(x, y)
	textOp.ColorScale.ScaleWithColor(color.White)
	textv2.Draw(screen, str, *resfont.GFace16, textOp)
}
