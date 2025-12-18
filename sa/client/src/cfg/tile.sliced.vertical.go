package cfg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

// TileVerticalSliced 垂直切片瓦片信息 (用于大型建筑的Y-Sorting)
type TileVerticalSliced struct {
	SlicePixel int              // 切片宽度 (像素)
	Slices     []*VerticalSlice // 切片列表
}

// VerticalSlice 单个垂直切片
type VerticalSlice struct {
	Image        *ebitenv2.Image // 切片图像 (高度=BottomPixelY+1)
	OffsetX      int             // X偏移量 (相对于原tile左边)
	Width        int             // 切片宽度
	BottomPixelY int             // 该切片最底部有效像素的Y坐标 (用于计算WY)
}

// NewTileVerticalSliced 创建垂直切片瓦片
func NewTileVerticalSliced() *TileVerticalSliced {
	return &TileVerticalSliced{
		Slices: make([]*VerticalSlice, 0),
	}
}

// Split 按像素宽度垂直切分瓦片
func (p *TileVerticalSliced) Split(fullImage *ebitenv2.Image, tileset *TiledTileset) {
	p.SlicePixel = tileset.VerticalSlicePixel
	totalWidth := tileset.TileWidth
	totalHeight := tileset.TileHeight

	// 按 slicePixel 宽度切分
	for x := 0; x < totalWidth; x += p.SlicePixel {
		sliceWidth := p.SlicePixel
		if x+sliceWidth > totalWidth {
			sliceWidth = totalWidth - x // 最后一片可能不足 slicePixel
		}

		// 扫描该切片，找最底部有像素的行
		bottomY := p.scanBottomPixel(fullImage, x, sliceWidth, totalHeight)

		if bottomY < 0 { // 该切片无有效像素，跳过
			continue
		}

		// 裁剪: 从 (x, 0) 到 (x+sliceWidth, bottomY+1)
		sliceRect := image.Rect(x, 0, x+sliceWidth, bottomY+1)
		sliceImage := fullImage.SubImage(sliceRect).(*ebitenv2.Image)

		slice := &VerticalSlice{
			Image:        sliceImage,
			OffsetX:      x,
			Width:        sliceWidth,
			BottomPixelY: bottomY,
		}
		p.Slices = append(p.Slices, slice)
	}
}

// scanBottomPixel 扫描指定区域，找到最底部有像素(非透明)的行
// 返回该行的Y坐标，如果没有有效像素则返回-1
func (p *TileVerticalSliced) scanBottomPixel(img *ebitenv2.Image, startX, width, height int) int {
	// 从底部向上扫描
	for y := height - 1; y >= 0; y-- {
		for x := startX; x < startX+width; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 { // 非透明像素
				return y
			}
		}
	}
	return -1 // 无有效像素
}
