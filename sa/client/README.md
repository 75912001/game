# SA

# todo
    tiled 编辑地图
    创建一个建筑物. 实现底部-不可到达区域,顶部-可到达区域
    实现遮挡效果

    在这个程序上. 我需要增加回合制战斗. 就是在地图上移动.
    随机触发遇敌,并在不离开地图的情况下. 切换到战斗场景. 战斗结束之后.
    再切换回刚才遇敌的场景. 那么应该如何设计. 不要修改我的代码. 而是给我设计方案

# 提示词
## 背景:
    现在使用 tiled 工具来制作地图. 现在需要地面图块 tile, 包含各种草地.
## 基础要求:
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
## 任务:
    制作规格为每个图块64*64
    图块之间无缝衔接
    制作16种草地,4*4排列
    最终出图为*.png

# 素材说明
## 
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

# 不可到达区域
## 实现不可到达区域的几种方案对比

|   方案   |            A.纯图块判断             |       B.对象层       |          C.图块集碰撞          |
|:------:|:------------------------------:|:-----------------:|:-------------------------:|
|   原理   |       检查坐标(tx,ty)处是否有图ID       |   手动在地图上画矩形/多边形   | 在图块集(Tileset)编辑器里给图块画好碰撞盒 |
|  制作效率  |           极高(画图即有碰撞)           |  低(每个房子都要重画一遍碰撞)  |       高(定义一次，处处生效)        |
|  运行效率  |          极高(O(1)数组查询)          |   中(需遍历/四叉树优化)    |     中(加载时生成数据，运行时同B)      |
|   精度   |           差(只能整格碰撞)            |     高(可自定义形状)     |         高(可自定义形状)         |
|   问题   |  角色会被空气墙挡住(因为图块是方形的,但树/柱子是圆的)  |      累,容易漏画       |          需要预处理加载          |
|  是否可用  |               ✅                |         ✅         |           ❌未实现            |


## 使用:
#### A 纯图块判断-图块层
    使用 tiled 的 图块层,增加自定义属性 blocked(bool), 即该图块层的所有图块,都为阻挡区域.
    如: Builing.1 图块层, blocked=true, 即,建筑物的最底层图块,都为阻挡区域.
#### B 对象层
    使用 tiled 的 对象层,手动在地图上画矩形, 并增加自定义属性 blocked(bool)
#### C 图块集碰撞
    使用 tiled 的 图块集编辑器, 给图块画好碰撞盒.
    这样, 每次使用该图块时, 都会自动带上碰撞盒.
    如: 树木图块, 柱子图块等.

#### 其他
    1. 使用透明图块, 设置为blocked区域
