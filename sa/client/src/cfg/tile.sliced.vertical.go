package cfg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

// TileSlicedVertical 瓦片-垂直切片 (用于大型建筑的Y-Sorting)
type TileSlicedVertical struct {
	SlicePixel int                 // 切片宽度 (像素)
	Slices     []*TileVerticalCell // 切片列表
}

// TileVerticalCell 单个垂直切片
type TileVerticalCell struct {
	Image        *ebitenv2.Image // 切片图像
	OffsetX      int             // X偏移量 (相对于原tile左边)
	Width        int             // 切片宽度
	BottomPixelY int             // 该切片最底部有效像素的Y坐标
}

// NewTileSlicedVertical 创建 瓦片-垂直切片
func NewTileSlicedVertical() *TileSlicedVertical {
	return &TileSlicedVertical{
		Slices: make([]*TileVerticalCell, 0),
	}
}

// Split 按像素宽度垂直切分瓦片
func (p *TileSlicedVertical) Split(fullImage *ebitenv2.Image, tileset *TiledTileset) {
	p.SlicePixel = tileset.VerticalSlicePixel
	totalWidth := tileset.TileWidth
	totalHeight := tileset.TileHeight

	// 按 slicePixel 宽度切分
	for offsetX := 0; offsetX < totalWidth; offsetX += p.SlicePixel {
		sliceWidth := p.SlicePixel
		if totalWidth < offsetX+sliceWidth { // 要切的宽度超出总宽度, 调整为剩余宽度
			sliceWidth = totalWidth - offsetX // 最后一片可能不足 slicePixel
		}

		// 扫描该切片，找最底部有像素的行
		bottomY := p.scanBottomPixel(fullImage, offsetX, sliceWidth, totalHeight)

		if bottomY < 0 { // 该切片无有效像素，跳过, 图像中间有一列透明像素时会出现这种情况, ⛔ 这里不应该出现
			panic("TileSlicedVertical.Split: slice has no valid pixels")
		}
		// 裁剪
		sliceRect := image.Rect(offsetX, 0, offsetX+sliceWidth, bottomY+1)
		sliceImage := fullImage.SubImage(sliceRect).(*ebitenv2.Image)

		slice := &TileVerticalCell{
			Image:        sliceImage,
			OffsetX:      offsetX,
			Width:        sliceWidth,
			BottomPixelY: bottomY,
		}
		p.Slices = append(p.Slices, slice)
	}
}

// scanBottomPixel 扫描指定区域，找到最底部有像素(非透明)的行-y
// 返回该行的Y坐标，如果没有有效像素则返回-1
func (p *TileSlicedVertical) scanBottomPixel(img *ebitenv2.Image, startX, width, height int) int {
	// 从底部向上扫描
	for y := height - 1; 0 <= y; y-- {
		for x := startX; x < startX+width; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if 0 < a { // 非透明像素
				return y
			}
		}
	}
	return -1 // 无有效像素
}
