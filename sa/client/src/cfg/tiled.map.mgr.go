package cfg

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	ct "saClient/src/coordinatetransform"
	"strconv"
	"strings"

	"saClient/src/common"
	"saClient/src/proto"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pkg/errors"
)

var GTiledMapMgr = newTiledMapMgr()

// TiledMapMgr Tiled 地图资源管理器
type TiledMapMgr struct {
	Maps *xmap.MapMgr[common.AssetID, *TiledMap]
}

func newTiledMapMgr() *TiledMapMgr {
	return &TiledMapMgr{
		Maps: xmap.NewMapMgr[common.AssetID, *TiledMap](),
	}
}

// Load 加载 Tiled 地图资源
func (p *TiledMapMgr) Load() error {
	tiledMapDir := common.AppResTiledMapDir
	// 检查目录是否存在
	if _, err := os.Stat(tiledMapDir); os.IsNotExist(err) {
		return errors.WithMessagef(err, "tiledMapDir %s does not exist", tiledMapDir)
	}

	mapFiles, err := os.ReadDir(tiledMapDir)
	if err != nil {
		return errors.WithMessagef(err, "读取 Tiled 地图目录失败: %v %v", tiledMapDir, xruntime.Location())
	}

	for _, mapFile := range mapFiles {
		if mapFile.IsDir() {
			continue
		}
		fileName := mapFile.Name()
		ext := filepath.Ext(fileName)
		if ext != ".tmx" { // 不是 .tmx 格式
			continue
		}
		// 解析文件名获取地图ID: map.${mapID}.tmx
		baseName := fileName[:len(fileName)-len(ext)]
		if len(baseName) < 5 || baseName[:4] != "map." { // 不符合命名规范 使用 map. 前缀
			continue
		}
		mapIDStr := baseName[4:]
		mapID64, err := strconv.ParseUint(mapIDStr, 10, 32)
		if err != nil {
			return errors.WithMessagef(err, "解析 Tiled 地图文件名 %s 为ID失败: %v", fileName, xruntime.Location())
		}
		mapID := common.AssetID(mapID64)
		// 检查ID范围
		if mapID < common.AssetID(proto.AssetIDRange_AssetIDRange_Map_Start) || common.AssetID(proto.AssetIDRange_AssetIDRange_Map_End) < mapID {
			return fmt.Errorf("Tiled 地图ID超出范围: %d %v", mapID, xruntime.Location())
		}
		// 加载地图
		tiledMap, err := p.loadTiledMap(mapID, filepath.Join(tiledMapDir, fileName))
		if err != nil {
			return errors.WithMessagef(err, "加载 Tiled 地图失败: %v %v", fileName, xruntime.Location())
		}
		ok := p.Maps.AddIfNotExist(mapID, tiledMap)
		if !ok {
			return fmt.Errorf("添加 Tiled 地图资源失败,地图已存在: %v %v", mapID, xruntime.Location())
		}
	}
	return nil
}

// loadTiledMap 加载单个 Tiled 地图 (TMX 格式)
func (p *TiledMapMgr) loadTiledMap(mapID common.AssetID, tmxPath string) (*TiledMap, error) {
	data, err := os.ReadFile(tmxPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "读取 TMX 文件失败: %v %v", tmxPath, xruntime.Location())
	}

	var mapXML tmxMap
	if err := xml.Unmarshal(data, &mapXML); err != nil {
		return nil, errors.WithMessagef(err, "解析 TMX XML 失败: %v %v", tmxPath, xruntime.Location())
	}
	if mapXML.Orientation != "isometric" {
		return nil, fmt.Errorf("仅支持等距(isometric)地图,当前地图 %v 类型: %s %v", mapID, mapXML.Orientation, xruntime.Location())
	}
	if mapXML.RenderOrder != "right-down" {
		return nil, fmt.Errorf("仅支持右下(right-down)渲染顺序,当前地图 %v 渲染顺序: %s %v", mapID, mapXML.RenderOrder, xruntime.Location())
	}
	if mapXML.Infinite != nil && *mapXML.Infinite != 0 {
		return nil, fmt.Errorf("不支持无限(infinite)地图,当前地图 %v infinite 属性: %v %v", mapID, mapXML.Infinite, xruntime.Location())
	}
	if mapXML.TileWidth/2 != mapXML.TileHeight {
		return nil, fmt.Errorf("仅支持等宽高比 2:1 的等距(isometric)地图,当前地图 %v 瓦片宽高: %d x %d %v", mapID, mapXML.TileWidth, mapXML.TileHeight, xruntime.Location())
	}

	tiledMap := &TiledMap{
		ID:         mapID,
		Width:      mapXML.Width,
		Height:     mapXML.Height,
		TileWidth:  mapXML.TileWidth,
		TileHeight: mapXML.TileHeight,
	}

	if mapXML.Properties != nil { // 解析地图属性
		for _, prop := range mapXML.Properties.Properties {
			switch prop.Name {
			case TiledMapTag_BgmFilePath: // 背景音乐文件路径
				tiledMap.BackgroundMusicFilePath = prop.Value
			default: // 未知属性 忽略
			}
		}
	}

	// 加载 tilesets
	tmxDir := filepath.Dir(tmxPath)
	for _, tsRef := range mapXML.Tilesets {
		tileset, err := p.loadTileset(tmxDir, tsRef)
		if err != nil {
			return nil, err
		}
		tiledMap.Tilesets = append(tiledMap.Tilesets, tileset) // todo menglc 这里不同的 mapID, 相同的 tileset 会重复加载,可以考虑缓存优化
	}

	// 加载 tile layers
	for _, layerXML := range mapXML.Layers {
		layer := &TiledLayer{
			ID:      layerXML.ID,
			Type:    TiledLayerType_TileLayer,
			Visible: layerXML.Visible == nil || *layerXML.Visible != 0,
			Opacity: layerXML.Opacity,
			Width:   layerXML.Width,
			Height:  layerXML.Height,
		}
		if layer.Opacity == 0 {
			layer.Opacity = 1.0
		}
		if layerXML.Data.Encoding == "csv" || layerXML.Data.Encoding == "" {
			// 解析非 infinite map 的 data
			layer.Data, err = parseCSVData(layerXML.Data.Content)
			if err != nil {
				return nil, errors.WithMessagef(err, "解析 %v 图层 %d 的 CSV 数据失败: %v %v", mapID, layer.ID, tmxPath, xruntime.Location())
			}
		}

		tiledMap.Layers = append(tiledMap.Layers, layer)
	}

	// 加载 object groups
	for _, ogXML := range mapXML.ObjectGroups {
		layer := &TiledLayer{
			ID:      ogXML.ID,
			Type:    TiledLayerType_ObjectLayer,
			Visible: ogXML.Visible == nil || *ogXML.Visible != 0,
			Opacity: 1.0,
		}

		for _, objXML := range ogXML.Objects {
			if !xutil.Float32Equal(objXML.Rotation, 0) {
				return nil, fmt.Errorf("为了提高效能, 不支持旋转对象, 当前地图 %v 对象 %d 旋转角度: %f %v", mapID, objXML.ID, objXML.Rotation, xruntime.Location())
			}
			switch objXML.Type {
			case TiledObjectType_Portal: // 传送点对象
			default:
				return nil, fmt.Errorf("不支持的对象类型, 当前地图 %v 对象 %d 类型: %s %v", mapID, objXML.ID, objXML.Type, xruntime.Location())
			}
			obj := &TiledObject{
				ID:      objXML.ID,
				Type:    objXML.Type,
				X:       objXML.X,
				Y:       objXML.Y,
				Width:   objXML.Width,
				Height:  objXML.Height,
				Visible: objXML.Visible == nil || *objXML.Visible != 0,
			}

			if objXML.Polygon != nil { // 多边形对象
				return nil, fmt.Errorf("为了提高效能, 不支持多边形对象, 当前地图 %v 对象 %d 类型: polygon %v", mapID, objXML.ID, xruntime.Location())
			}

			// 解析对象属性
			if objXML.Properties != nil {
				for _, prop := range objXML.Properties.Properties {
					switch prop.Name {
					case TiledObjectTag_Collision:
						obj.Collision = prop.Value == "true"
					case TiledObjectTag_TargetPortal:
						v, err := strconv.ParseUint(prop.Value, 10, 32)
						if err != nil {
							return nil, errors.WithMessagef(err, "解析对象 %d 的目标传送点失败: %v %v", obj.ID, tmxPath, xruntime.Location())
						}
						obj.TargetPortal = common.PortalID(v)

					}
				}
			}

			layer.Objects = append(layer.Objects, obj)
		}

		tiledMap.Layers = append(tiledMap.Layers, layer)
	}

	return tiledMap, nil
}

// loadTileset 加载 tileset (TSX 格式)
func (p *TiledMapMgr) loadTileset(tmxDir string, tsRef tmxTilesetRef) (*TiledTileset, error) {
	tsxPath := filepath.Join(tmxDir, tsRef.Source)
	tsxPath = filepath.Clean(tsxPath)

	data, err := os.ReadFile(tsxPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "读取 TSX 文件失败: %v %v", tsxPath, xruntime.Location())
	}

	var tsxXML tsxTileset
	if err := xml.Unmarshal(data, &tsxXML); err != nil {
		return nil, errors.WithMessagef(err, "解析 TSX XML 失败: %v %v", tsxPath, xruntime.Location())
	}

	tileset := &TiledTileset{
		FirstGID:    tsRef.FirstGID,
		Name:        tsxXML.Name,
		TileWidth:   tsxXML.TileWidth,
		TileHeight:  tsxXML.TileHeight,
		TileCount:   tsxXML.TileCount,
		Columns:     tsxXML.Columns,
		ImageWidth:  tsxXML.Image.Width,
		ImageHeight: tsxXML.Image.Height,
	}

	// 加载图片
	tsxDir := filepath.Dir(tsxPath)
	imgPath := filepath.Join(tsxDir, tsxXML.Image.Source)
	imgPath = filepath.Clean(imgPath)

	img, _, err := ebitenv2ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "加载 tileset 图片失败: %v %v", imgPath, xruntime.Location())
	}
	tileset.Image = img

	return tileset, nil
}

// parseCSVData 解析 CSV 格式的 tile 数据
func parseCSVData(content string) ([]int, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, nil
	}

	var result []int
	parts := strings.Split(content, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, errors.WithMessagef(err, "解析 CSV tile 数据失败: %v %v", err, xruntime.Location())
		}
		result = append(result, val)
	}
	return result, nil
}

// Check 检查 Tiled 地图资源
func (p *TiledMapMgr) Check() error {
	var err error
	p.Maps.Foreach(
		func(mapID common.AssetID, tiledMap *TiledMap) (isContinue bool) {
			exist := GMapMgr.Maps.IsExist(mapID)
			if !exist { // 检查地图是否合法
				err = fmt.Errorf("Tiled 地图资源 %v 未定义 %v", mapID, xruntime.Location())
				return false
			}
			for _, layer := range tiledMap.Layers {
				switch layer.Type {
				case TiledLayerType_TileLayer: // 瓦片图层
				case TiledLayerType_ObjectLayer: // 对象图层
					for _, obj := range layer.Objects {
						switch obj.Type {
						case TiledObjectType_Portal: // 传送点对象
							if obj.TargetPortal == 0 { // 必须设置目标传送点
								err = fmt.Errorf("Tiled 地图资源 %v 中对象 %d 的目标传送点未设置 %v", mapID, obj.ID, xruntime.Location())
								return false
							}
							exist = GPortalMgr.Points.IsExist(obj.TargetPortal)
							if !exist { // 检查目标传送点是否存在
								err = fmt.Errorf("Tiled 地图资源 %v 中对象 %d 的目标传送点 %v 不存在 %v", mapID, obj.ID, obj.TargetPortal, xruntime.Location())
								return false
							}
						}
					}
				}
			}
			return true
		},
	)
	return err
}

// Assemble 装配 Tiled 地图资源
func (p *TiledMapMgr) Assemble() error {
	p.Maps.Foreach(
		func(mapID common.AssetID, tiledMap *TiledMap) (isContinue bool) {
			// 等距地图像素尺寸 (包含完整的菱形内容区域)
			// +TileHeight 是为了包含第一行 tile 上方和最后一行 tile 下方的菱形边缘
			tiledMap.PixelW = (tiledMap.Width + tiledMap.Height) * (tiledMap.TileWidth / 2)
			tiledMap.PixelH = (tiledMap.Width + tiledMap.Height) * (tiledMap.TileHeight / 2)
			tiledMap.IsometricCT = ct.NewIsometric(tiledMap.Width, tiledMap.Height, tiledMap.TileWidth, tiledMap.TileHeight)
			return true
		},
	)

	return nil
}
