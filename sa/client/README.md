# SA

# 素材说明
`
2.5D mmorpg 游戏素材
伪3d(Fake 3D)
视角:等轴测视角 45度 (Isometric view)
风格:角色二次元/卡通,非真实感渲染(NPR),赛璐璐渲染(Cel-shaded)
矢量艺术(Vector art)
线条清晰
加粗轮廓(Bold outlines)
平涂色(Flat color)
高饱和度(High saturation)
色彩鲜艳
整体色调明亮,温暖
画面干净
纹理简洁
细节:极简纹理(Minimal texture)
高质量
营造出轻松,愉快的童话冒险氛围
背景透明
`


背景:
    现在使用 tiled 工具来制作地图. 现在需要地面图块 tile, 包含各种草地.
基础要求:
    2.5D mmorpg 游戏素材
    伪3d(Fake 3D)
    视角:等轴测视角 45度 (Isometric view)
    风格:角色二次元/卡通,非真实感渲染(NPR),赛璐璐渲染(Cel-shaded)
    矢量艺术(Vector art)
    线条清晰
    加粗轮廓(Bold outlines)
    平涂色(Flat color)
    高饱和度(High saturation)
    色彩鲜艳
    整体色调明亮,温暖
    画面干净
    纹理简洁
    细节:极简纹理(Minimal texture)
    高质量
    营造出轻松,愉快的童话冒险氛围
    背景透明
任务:
    制作规格为每个图块64*64
    图块之间无缝衔接
    制作16种草地,4*4排列
    最终出图为*.png





# todo

// drawData 绘制 tile 数据
func (p *TiledMap) drawData(screen *ebitenv2.Image, cam *camera.Camera, data []int, startX, startY, width, height int) {
研究逻辑, 优化计算, 让相关的瓦片挂载在地图上, 不用每次绘制的时候都放在内存中. 每次加载地图.都将该地图资源全量加载, 离开地图时,卸载地图所占资源

tiled 编辑地图
    制作2张地图
        1. 新手村
            2个传送点
                1. 野外.位置1
                2. 野外.位置2
        2. 野外
            2个传送点
                1. 新手村.位置1
                2. 新手村.位置2
开发地图跳转功能

# 相关工具
## 地图编辑器 Tiled
    https://www.mapeditor.org/
    正常 orthogonal
    ✅45度 isometric
    等角(交错) staggered
    

## 解析 Tiled 地图(TMX)的库 go-tiled
    github.com/lafriks/go-tiled
    Go library to parse Tiled map editor file format (TMX) and render map to image
## 精灵图集打包工具 TexturePacker
    https://www.codeandweb.com/texturepacker
## PhotoScape X (图像编辑软件)
    http://x.photoscape.org/
## 移除图像背景 Remove Image Backgrounds
    https://pixian.ai/
## 精灵播放器 Sprite Reel
    https://spritereel.com/