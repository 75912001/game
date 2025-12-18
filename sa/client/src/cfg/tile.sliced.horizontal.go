package cfg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

// TileSlicedHorizontal 瓦片-水平切片 (用于标记需要拆分的瓦片类型)
type TileSlicedHorizontal struct {
	CanopyImage  *ebitenv2.Image // 上部，绘制在 Overhead 层
	CanopyHeight int             // 上部高度
	TrunkImage   *ebitenv2.Image // 下部，参与 Y-Sorting
	TrunkHeight  int             // 下部高度
}

func NewTileSliceHorizontal() *TileSlicedHorizontal {
	return &TileSlicedHorizontal{}
}

// Split 拆分
func (p *TileSlicedHorizontal) Split(fullImage *ebitenv2.Image, tileset *TiledTileset) {
	p.CanopyHeight = tileset.HorizontalSlicePixel
	p.TrunkHeight = tileset.TileHeight - p.CanopyHeight

	canopyRect := image.Rect(0, 0, tileset.TileWidth, p.CanopyHeight)
	p.CanopyImage = fullImage.SubImage(canopyRect).(*ebitenv2.Image)

	trunkRect := image.Rect(0, p.CanopyHeight, tileset.TileWidth, tileset.TileHeight)
	p.TrunkImage = fullImage.SubImage(trunkRect).(*ebitenv2.Image)
}
