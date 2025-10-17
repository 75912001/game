package font

import (
	apccommon "airplaneClient/src/common"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
	"os"
	"path/filepath"
)

var GFace32 *text.Face
var GFace20 *text.Face
var GFace16 *text.Face

var GFaceButton *text.Face // 按钮

func init() {
	// 加载中文字体
	var fontSource *text.GoTextFaceSource
	var err error
	// 尝试从资源目录加载中文字体
	chineseFontPath := filepath.Join(apccommon.AppResourcesFontsDir, "fangzheng_kaiti.ttf")
	fontData, err := os.ReadFile(chineseFontPath)
	if err != nil {
		log.Fatal(err)
	}
	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}
	{
		var face text.Face = &text.GoTextFace{
			Source: fontSource,
			Size:   32, // 字体大小
		}
		GFace32 = &face
	}
	{ // 创建字体 Face
		// 将 GoTextFace 转换为 *text.Face (指向接口的指针)
		var face text.Face = &text.GoTextFace{
			Source: fontSource,
			Size:   20, // 字体大小
		}
		GFace20 = &face
	}
	{
		var face text.Face = &text.GoTextFace{
			Source: fontSource,
			Size:   16, // 字体大小
		}
		GFace16 = &face
	}
	GFaceButton = GFace32
}
