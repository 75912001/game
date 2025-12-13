package cfg

import (
	ct "saClient/src/coordinatetransform"

	"saClient/src/common"
)

type TiledLayerType uint32 // Tiled图层类型

const (
	TiledLayerType_Unknown     TiledLayerType = 0 // 未知图层
	TiledLayerType_TileLayer   TiledLayerType = 1 // 瓦片图层
	TiledLayerType_ObjectLayer TiledLayerType = 2 // 对象图层
)

// TiledMap 属性标签
const TiledMapTag_BgmFilePath = "backgroundMusicFilePath" // 背景音乐文件路径-标签

// TiledMap Tiled 地图资源
type TiledMap struct {
	ID                      common.AssetID  // 地图ID
	Width                   int             // 地图宽度（tile 数）
	Height                  int             // 地图高度（tile 数）
	TileWidth               int             // tile 宽度（像素）
	TileHeight              int             // tile 高度（像素）
	Layers                  []*TiledLayer   // 图层列表
	Tilesets                []*TiledTileset // tileset 列表
	BackgroundMusicFilePath string          // 背景音乐文件路径

	PixelW      int           // 地图像素宽度
	PixelH      int           // 地图像素高度
	IsometricCT *ct.Isometric // 坐标转换器
}

// TiledLayer Tiled 图层
type TiledLayer struct {
	ID      int            // 图层ID
	Type    TiledLayerType // 图层类型
	Visible bool           // 是否可见
	Opacity float32        // 透明度
	X       int            // 图层X偏移
	Y       int            // 图层Y偏移
	Width   int            // 图层宽度
	Height  int            // 图层高度
	Data    []int          // tile-数据-用于有限地图
	Objects []*TiledObject // 对象列表
}
