package font

import (
	"bytes"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"os"
	"path/filepath"
	"saClient/src/common"
)

var GFace32 *textv2.Face
var GFace20 *textv2.Face
var GFace16 *textv2.Face

var GFaceButton *textv2.Face // 按钮

func init() {
	// 加载中文字体
	var fontSource *textv2.GoTextFaceSource
	var err error
	// 尝试从资源目录加载中文字体
	chineseFontPath := filepath.Join(common.AppResFontDir, "fangzheng_kaiti.ttf")
	fontData, err := os.ReadFile(chineseFontPath)
	if err != nil {
		panic(err)
	}
	fontSource, err = textv2.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic(err)
	}
	{
		var face textv2.Face = &textv2.GoTextFace{
			Source: fontSource,
			Size:   32, // 字体大小
		}
		GFace32 = &face
	}
	{ // 创建字体 Face
		// 将 GoTextFace 转换为 *text.Face (指向接口的指针)
		var face textv2.Face = &textv2.GoTextFace{
			Source: fontSource,
			Size:   20, // 字体大小
		}
		GFace20 = &face
	}
	{
		var face textv2.Face = &textv2.GoTextFace{
			Source: fontSource,
			Size:   16, // 字体大小
		}
		GFace16 = &face
	}
	GFaceButton = GFace32
}
