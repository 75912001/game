package tiled

import (
	"image"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"saClient/src/common"
)

// TiledMap Tiled 地图资源
type TiledMap struct {
	ID                      common.AssetID  // 地图ID
	Width                   int             // 地图宽度（tile 数）
	Height                  int             // 地图高度（tile 数）
	TileWidth               int             // tile 宽度（像素）
	TileHeight              int             // tile 高度（像素）
	Orientation             string          // 地图方向: isometric (等距)
	Layers                  []*TiledLayer   // 图层列表
	Tilesets                []*TiledTileset // tileset 列表
	BackgroundMusicFilePath string          // 背景音乐文件路径
}

// TiledLayer Tiled 图层
type TiledLayer struct {
	ID      int            // 图层ID
	Type    LayerType      // 图层类型
	Visible bool           // 是否可见
	Opacity float64        // 透明度
	X       int            // 图层X偏移
	Y       int            // 图层Y偏移
	Width   int            // 图层宽度
	Height  int            // 图层高度
	Data    []int          // tile 数据（用于有限地图）
	Objects []*TiledObject // 对象列表（用于 objectgroup）
}

// TiledObject Tiled 对象（用于碰撞等）
type TiledObject struct {
	ID        int           // 对象ID
	Type      string        // 对象类型
	X         float64       // 对象X坐标
	Y         float64       // 对象Y坐标
	Width     float64       // 对象宽度
	Height    float64       // 对象高度
	Rotation  float64       // 旋转角度
	Visible   bool          // 是否可见
	Polygon   []*TiledPoint // 多边形顶点（相对坐标）
	Collision bool          // 是否碰撞
}

// TiledPoint Tiled 点
type TiledPoint struct {
	X float64
	Y float64
}

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

// GetPixelWidth 获取地图像素宽度
func (m *TiledMap) GetPixelWidth() int {
	switch m.Orientation {
	case "isometric":
		// 等距地图(菱形):宽度 = (Width + Height) * TileWidth / 2
		return (m.Width + m.Height) * (m.TileWidth / 2)
	default:
		return m.Width * m.TileWidth
	}
}

// GetPixelHeight 获取地图像素高度
func (m *TiledMap) GetPixelHeight() int {
	switch m.Orientation {
	case "isometric":
		// 等距地图(菱形):高度 = (Width + Height) * TileHeight / 2
		return (m.Width + m.Height) * (m.TileHeight / 2)
	default:
		return m.Height * m.TileHeight
	}
}
