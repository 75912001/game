package cfg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

// TileUmbrellaShaped 伞形瓦片信息-上下拆分 (用于标记需要拆分的瓦片类型)
type TileUmbrellaShaped struct {
	CanopyImage  *ebitenv2.Image // 冠部图像 (上部分，绘制在 Overhead 层)
	CanopyHeight int             // 冠部高度
	TrunkImage   *ebitenv2.Image // 干部图像 (下部分，参与 Y-Sorting)
	TrunkHeight  int             // 干部高度
}

func NewTileUmbrellaShaped() *TileUmbrellaShaped {
	return &TileUmbrellaShaped{}
}

// Split 拆分伞形瓦片
func (p *TileUmbrellaShaped) Split(fullImage *ebitenv2.Image, tileset *TiledTileset) {
	// 计算拆分高度
	totalHeight := tileset.TileHeight
	canopyHeight := int(float32(totalHeight) * tileset.OverheadRatio)
	trunkHeight := totalHeight - canopyHeight
	// 确保最小高度
	if canopyHeight < 1 || trunkHeight < 1 {
		return // 不拆分
	}
	p.CanopyHeight = canopyHeight
	p.TrunkHeight = trunkHeight

	canopyRect := image.Rect(0, 0, tileset.TileWidth, canopyHeight)
	p.CanopyImage = fullImage.SubImage(canopyRect).(*ebitenv2.Image)

	trunkRect := image.Rect(0, canopyHeight, tileset.TileWidth, totalHeight)
	p.TrunkImage = fullImage.SubImage(trunkRect).(*ebitenv2.Image)
}
