package cfg

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// TiledTileBlocked tile 阻挡定义 (相对于 tile 左上角的偏移)
type TiledTileBlocked struct {
	X      float32 // 相对 tile 左上角的 X 偏移
	Y      float32 // 相对 tile 左上角的 Y 偏移
	Width  float32 // 阻挡体宽度
	Height float32 // 阻挡体高度
}

// TiledTileset Tiled-瓦片集
type TiledTileset struct {
	FirstGID    int             // 起始 GID
	TileWidth   int             // tile 宽度
	TileHeight  int             // tile 高度
	TileCount   int             // tile 总数
	Columns     int             // 列数
	Image       *ebitenv2.Image // tileset 图片
	ImageWidth  int             // 图片宽度
	ImageHeight int             // 图片高度

	HorizontalSlicePixel int // 水平切分像素值 (只切一次,上下切分)
	VerticalSlicePixel   int // 垂直切分像素值 (切多次, 64=每64像素垂直切分一次)

	TileBlocked map[int][]*TiledTileBlocked // tile 阻挡体 (key: tile 本地 ID)
}
