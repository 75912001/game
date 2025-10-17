package ui

import (
	apcfont "airplaneClient/src/resources/font"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

func Printf(screen *ebiten.Image, x, y float64, format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(x, y)
	textOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, str, *apcfont.GFace16, textOp)
}
