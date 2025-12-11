package tiled

type LayerType uint32 // Tiled图层类型

const (
	LayerType_Unknown     LayerType = 0 // 未知图层
	LayerType_TileLayer   LayerType = 1 // 瓦片图层
	LayerType_ObjectLayer LayerType = 2 // 对象图层
)

const MapBgmFilePathTag = "backgroundMusicFilePath" // 背景音乐文件路径-标签
const ObjectCollisionTag = "collision"              // 碰撞对象-标签
