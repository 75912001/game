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


tiled 编辑地图


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

在这个程序上. 我需要增加回合制战斗. 就是在地图上移动.
随机触发遇敌,并在不离开地图的情况下. 切换到战斗场景. 战斗结束之后.
再切换回刚才遇敌的场景. 那么应该如何设计. 不要修改我的代码. 而是给我设计方案 