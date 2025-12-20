package cfg

import (
	"encoding/xml"
)

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
	ID         int            `xml:"id,attr"`
	Class      string         `xml:"class,attr"`
	Width      int            `xml:"width,attr"`
	Height     int            `xml:"height,attr"`
	Visible    *int           `xml:"visible,attr"`
	Opacity    float32        `xml:"opacity,attr"`
	Data       tmxData        `xml:"data"`
	Properties *tmxProperties `xml:"properties"`
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
	Class   string      `xml:"class,attr"`
	Visible *int        `xml:"visible,attr"`
	Objects []tmxObject `xml:"object"`
}

type tmxObject struct {
	ID         int             `xml:"id,attr"`
	Type       TiledObjectType `xml:"type,attr"`
	X          float32         `xml:"x,attr"`
	Y          float32         `xml:"y,attr"`
	Width      float32         `xml:"width,attr"`
	Height     float32         `xml:"height,attr"`
	Rotation   float32         `xml:"rotation,attr"`
	Visible    *int            `xml:"visible,attr"`
	Polygon    *tmxPolygon     `xml:"polygon"`
	Properties *tmxProperties  `xml:"properties"`
}

type tmxPolygon struct {
	Points string `xml:"points,attr"`
}

// TSX XML 结构定义
type tsxTileset struct {
	XMLName    xml.Name       `xml:"tileset"`
	TileWidth  int            `xml:"tilewidth,attr"`
	TileHeight int            `xml:"tileheight,attr"`
	TileCount  int            `xml:"tilecount,attr"`
	Columns    int            `xml:"columns,attr"`
	Image      tsxImage       `xml:"image"`
	Properties *tmxProperties `xml:"properties"` // tileset 属性 (如 overheadRatio)
	Tiles      []tsxTile      `xml:"tile"`       // tile 碰撞体定义
}

// tsxTile 单个 tile 定义(包含碰撞体)
type tsxTile struct {
	ID          int             `xml:"id,attr"`
	ObjectGroup *tmxObjectGroup `xml:"objectgroup"` // tile 碰撞体
}

type tsxImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}
