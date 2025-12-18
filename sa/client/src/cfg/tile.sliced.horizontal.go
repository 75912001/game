package cfg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

// TileSlicedHorizontal 瓦片-水平切片 (用于标记需要拆分的瓦片类型)
type TileSlicedHorizontal struct {
	CanopyImage  *ebitenv2.Image // 冠部图像 (上部分，绘制在 Overhead 层)
	CanopyHeight int             // 冠部高度
	TrunkImage   *ebitenv2.Image // 干部图像 (下部分，参与 Y-Sorting)
	TrunkHeight  int             // 干部高度
}

func NewTileSliceHorizontal() *TileSlicedHorizontal {
	return &TileSlicedHorizontal{}
}

// Split 拆分
func (p *TileSlicedHorizontal) Split(fullImage *ebitenv2.Image, tileset *TiledTileset) {
	p.CanopyHeight, p.TrunkHeight = CalculateSliceHorizontalSplitHeight(tileset.TileHeight, tileset.HorizontalSlicePixel) // = canopyHeight
	p.TrunkHeight = tileset.TileHeight - p.CanopyHeight

	canopyRect := image.Rect(0, 0, tileset.TileWidth, p.CanopyHeight)
	p.CanopyImage = fullImage.SubImage(canopyRect).(*ebitenv2.Image)

	trunkRect := image.Rect(0, p.CanopyHeight, tileset.TileWidth, tileset.TileHeight)
	p.TrunkImage = fullImage.SubImage(trunkRect).(*ebitenv2.Image)
}

// 计算拆分高度, 返回冠部和干部高度
func CalculateSliceHorizontalSplitHeight(totalHeight int, overheadRatio float32) (canopyHeight int, trunkHeight int) {
	canopyHeight = int(float32(totalHeight) * overheadRatio)
	trunkHeight = totalHeight - canopyHeight
	// 确保最小高度
	if canopyHeight < 1 || trunkHeight < 1 {
		return 0, 0 // 不拆分
	}
	return canopyHeight, trunkHeight
}
