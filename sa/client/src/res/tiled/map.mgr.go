package tiled

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	ebitenv2ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pkg/errors"
	"saClient/src/common"
	"saClient/src/proto"
)

var GMapMgr = newMapMgr()

// MapMgr Tiled 地图资源管理器
type MapMgr struct {
	Maps *xmap.MapMgr[common.AssetID, *TiledMap]
}

func newMapMgr() *MapMgr {
	return &MapMgr{
		Maps: xmap.NewMapMgr[common.AssetID, *TiledMap](),
	}
}

// TMX XML 结构定义
type tmxMap struct {
	Orientation  string           `xml:"orientation,attr"`
	RenderOrder  string           `xml:"renderorder,attr"`
	Width        int              `xml:"width,attr"`
	Height       int              `xml:"height,attr"`
	TileWidth    int              `xml:"tilewidth,attr"`
	TileHeight   int              `xml:"tileheight,attr"`
	Infinite     *int             `xml:"infinite,attr"`
	Properties   *tmxProperties   `xml:"properties"`
	Tilesets     []tmxTilesetRef  `xml:"tileset"`
	Layers       []tmxLayer       `xml:"layer"`
	ObjectGroups []tmxObjectGroup `xml:"objectgroup"`
}

// tmxProperties 地图属性集合
type tmxProperties struct {
	Properties []tmxProperty `xml:"property"`
}

// tmxProperty 单个属性
type tmxProperty struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type tmxTilesetRef struct {
	FirstGID int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr"`
}

type tmxLayer struct {
	ID      int     `xml:"id,attr"`
	Width   int     `xml:"width,attr"`
	Height  int     `xml:"height,attr"`
	Visible *int    `xml:"visible,attr"`
	Opacity float64 `xml:"opacity,attr"`
	Data    tmxData `xml:"data"`
}

type tmxData struct {
	Encoding string     `xml:"encoding,attr"`
	Chunks   []tmxChunk `xml:"chunk"`
	Content  string     `xml:",chardata"`
}

type tmxChunk struct {
	X       int    `xml:"x,attr"`
	Y       int    `xml:"y,attr"`
	Width   int    `xml:"width,attr"`
	Height  int    `xml:"height,attr"`
	Content string `xml:",chardata"`
}

type tmxObjectGroup struct {
	ID      int         `xml:"id,attr"`
	Visible *int        `xml:"visible,attr"`
	Objects []tmxObject `xml:"object"`
}

type tmxObject struct {
	ID         int            `xml:"id,attr"`
	Type       string         `xml:"type,attr"`
	X          float64        `xml:"x,attr"`
	Y          float64        `xml:"y,attr"`
	Width      float64        `xml:"width,attr"`
	Height     float64        `xml:"height,attr"`
	Rotation   float64        `xml:"rotation,attr"`
	Visible    *int           `xml:"visible,attr"`
	Polygon    *tmxPolygon    `xml:"polygon"`
	Properties *tmxProperties `xml:"properties"`
}

type tmxPolygon struct {
	Points string `xml:"points,attr"`
}

// TSX XML 结构定义
type tsxTileset struct {
	XMLName    xml.Name `xml:"tileset"`
	Name       string   `xml:"name,attr"`
	TileWidth  int      `xml:"tilewidth,attr"`
	TileHeight int      `xml:"tileheight,attr"`
	TileCount  int      `xml:"tilecount,attr"`
	Columns    int      `xml:"columns,attr"`
	Image      tsxImage `xml:"image"`
}

type tsxImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// Load 加载 Tiled 地图资源
func (p *MapMgr) Load() error {
	tiledMapDir := common.AppResTiledMapDir
	// 检查目录是否存在
	if _, err := os.Stat(tiledMapDir); os.IsNotExist(err) {
		return nil // 目录不存在则跳过加载-符合:无 Tiled 地图资源,则使用普通地图
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
func (p *MapMgr) loadTiledMap(mapID common.AssetID, tmxPath string) (*TiledMap, error) {
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
		ID:          mapID,
		Width:       mapXML.Width,
		Height:      mapXML.Height,
		TileWidth:   mapXML.TileWidth,
		TileHeight:  mapXML.TileHeight,
		Orientation: mapXML.Orientation,
	}

	if mapXML.Properties != nil { // 解析地图属性
		for _, prop := range mapXML.Properties.Properties {
			switch prop.Name {
			case MapBgmFilePathTag: // 背景音乐文件路径
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
			Type:    LayerType_TileLayer,
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
			Type:    LayerType_ObjectLayer,
			Visible: ogXML.Visible == nil || *ogXML.Visible != 0,
			Opacity: 1.0,
		}

		for _, objXML := range ogXML.Objects {
			obj := &TiledObject{
				ID:       objXML.ID,
				Type:     objXML.Type,
				X:        objXML.X,
				Y:        objXML.Y,
				Width:    objXML.Width,
				Height:   objXML.Height,
				Rotation: objXML.Rotation,
				Visible:  objXML.Visible == nil || *objXML.Visible != 0,
			}

			// 解析多边形
			if objXML.Polygon != nil {
				obj.Polygon = parsePolygonPoints(objXML.Polygon.Points)
			}

			// 解析对象属性
			if objXML.Properties != nil {
				for _, prop := range objXML.Properties.Properties {
					switch prop.Name {
					case ObjectCollisionTag:
						obj.Collision = prop.Value == "true"
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
func (p *MapMgr) loadTileset(tmxDir string, tsRef tmxTilesetRef) (*TiledTileset, error) {
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

// parsePolygonPoints 解析多边形点 "x1,y1 x2,y2 x3,y3"
func parsePolygonPoints(points string) []*TiledPoint {
	points = strings.TrimSpace(points)
	if points == "" {
		return nil
	}

	var result []*TiledPoint
	pairs := strings.Split(points, " ")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		coords := strings.Split(pair, ",")
		if len(coords) != 2 {
			continue
		}
		x, err1 := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
		y, err2 := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
		if err1 != nil || err2 != nil {
			continue
		}
		result = append(result, &TiledPoint{X: x, Y: y})
	}
	return result
}

// Check 检查 Tiled 地图资源
func (p *MapMgr) Check() error {
	// todo menglc 检查该资源是否-被地图配置引用
	return nil
}

// Assemble 装配 Tiled 地图资源
func (p *MapMgr) Assemble() error {
	return nil
}
