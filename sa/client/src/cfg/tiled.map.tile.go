package cfg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"image"
)

// TiledTileset Tiled-瓦片集
type TiledTileset struct {
	FirstGID    int             // 起始 GID
	Name        string          // tileset 名称
	TileWidth   int             // tile 宽度
	TileHeight  int             // tile 高度
	TileCount   int             // tile 总数
	Columns     int             // 列数
	Image       *ebitenv2.Image // tileset 图片
	ImageWidth  int             // 图片宽度
	ImageHeight int             // 图片高度
}

// GetTileImage 从 tileset 获取指定 tile 的子图
func (t *TiledTileset) GetTileImage(localID int) *ebitenv2.Image {
	if t.Image == nil || localID < 0 || localID >= t.TileCount {
		return nil
	}
	col := localID % t.Columns
	row := localID / t.Columns
	x := col * t.TileWidth
	y := row * t.TileHeight
	return t.Image.SubImage(image.Rect(x, y, x+t.TileWidth, y+t.TileHeight)).(*ebitenv2.Image)
}
