package ui

import (
	"fmt"
	imagecolor "image/color"
	resfont "saClient/src/res/font"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
)

func Printf(screen *ebitenv2.Image, x, y float32, format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	textOp := &textv2.DrawOptions{}
	textOp.GeoM.Translate(float64(x), float64(y))
	textOp.ColorScale.ScaleWithColor(imagecolor.White)
	textv2.Draw(screen, str, *resfont.GFace16, textOp)
}
