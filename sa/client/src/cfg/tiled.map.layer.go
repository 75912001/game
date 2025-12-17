package cfg

import (
	xmap "github.com/75912001/xlib/map"
)

type TiledLayerType uint32 // Tiled图层类型

const (
	TiledLayerType_Unknown   TiledLayerType = 0 // 未知图层
	TiledLayerType_Collision TiledLayerType = 1 // 碰撞图层

	TiledLayerType_Ground   TiledLayerType = 101 // 地面图层
	TiledLayerType_Building TiledLayerType = 111 // 建筑图层
	TiledLayerType_Objects  TiledLayerType = 121 // 物体图层
	TiledLayerType_Overhead TiledLayerType = 131 // 头顶图层
)

var TiledLayerTypeMap = xmap.NewMapMgr[TiledLayerType, string]()
var TiledLayerClassNameMap = xmap.NewMapMgr[string, TiledLayerType]()

func init() {
	TiledLayerTypeMap.Add(TiledLayerType_Collision, "Collision")
	TiledLayerTypeMap.Add(TiledLayerType_Ground, "Ground")
	TiledLayerTypeMap.Add(TiledLayerType_Building, "Building")
	TiledLayerTypeMap.Add(TiledLayerType_Objects, "Objects")
	TiledLayerTypeMap.Add(TiledLayerType_Overhead, "Overhead")

	TiledLayerClassNameMap.Add("Collision", TiledLayerType_Collision)
	TiledLayerClassNameMap.Add("Ground", TiledLayerType_Ground)
	TiledLayerClassNameMap.Add("Building", TiledLayerType_Building)
	TiledLayerClassNameMap.Add("Objects", TiledLayerType_Objects)
	TiledLayerClassNameMap.Add("Overhead", TiledLayerType_Overhead)
}

// TiledLayer Tiled 图层
type TiledLayer struct {
	ID           int            // 图层ID
	LayerType    TiledLayerType // 图层类型
	Visible      bool           // 是否可见
	Opacity      float32        // 透明度
	X            int            // 图层X偏移
	Y            int            // 图层Y偏移
	Width        int            // 图层宽度
	Height       int            // 图层高度
	Data         []int          // tile-数据-用于有限地图
	Objects      []*TiledObject // 对象列表
	BlockedLayer bool           // 是否为阻挡图层 (该图层所有非空 tile 都 阻挡)
}
