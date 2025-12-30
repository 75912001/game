# 基本要求
    深度思考 
    分析背景
    明确问题
    评审当前想法
    查阅可参考资料
    完成任务``
    给出详细的解决方案.
    注重逻辑清晰和步骤明确.
    不要直接给出代码实现, 而是给出如何实现的思路.
    注重性能和可维护性.
    如若有更好的方案, 请一并给出.并给出理由和对比.
    可以完全访问当前项目的代码和资源.
# 背景
    当前 arpg, 地图有生成怪物组的点.``````````
# 问题
    
# 当前想法
    以生成怪物组的点为圆心,画出生成怪物半径的圆(虚线)(红颜色)
    以生成怪物组的点为圆心,画出巡逻半径的圆(虚线)(橙颜色)
    以生成怪物组的点为圆心,画出追击半径的圆(虚线)(黄颜色)
# 可参考资料
# 任务
    给出方案.

# 解决方案 

## 方案概述
在 `scene.map.go` 中添加一个新的调试绘制方法 `drawSpawnPointDebug`，用于绘制刷怪点的各种半径圆。

## 具体实现步骤

### 1. 绘制虚线圆的辅助函数
由于 `ebiten/v2/vector` 没有直接提供虚线圆绘制功能，需要自己实现：
- **方式1（推荐）**：使用多段短弧线模拟虚线效果。将圆分成若干段，每隔一段绘制一小段弧线。
- **方式2**：直接使用 `vector.StrokeLine` 绘制多个短线段连成圆形轮廓，通过间隔控制虚线效果。

### 2. 数据来源
从 `ArpgEnemySpawnMgr` 获取所有刷怪点，每个刷怪点 (`ArpgEnemySpawnPoint`) 包含：
- `Object.WX`, `Object.WY` - 刷怪点世界坐标（圆心）
- `Object.SpawnRadius` - 生成半径（红色）
- `Object.PatrolRadius` - 巡逻半径（橙色）
- 追击半径 = `PatrolRadius * chaseRadiusMultiplier (2.0)` - 黄色

### 3. 坐标转换
将世界坐标 (WX, WY) 转换为屏幕坐标后绘制：
```
screenX, screenY = IsometricCT.W2S(WX, WY, camX, camY)
```

### 4. 绘制虚线圆的实现
```go
// drawDashedCircle 绘制虚线圆
// cx, cy: 圆心屏幕坐标
// radius: 半径（世界坐标单位，需要根据实际情况缩放）
// clr: 颜色
// dashCount: 虚线段数量（建议 32-64）
// dashRatio: 实线占比（0.5 表示实线和空白各占一半）
func drawDashedCircle(screen *ebitenv2.Image, cx, cy, radius float32, clr color.RGBA, dashCount int, dashRatio float32)
```

### 5. 集成到 DrawCollision
在 `Map.DrawCollision` 方法中调用 `drawSpawnPointDebug`：
```go
func (p *Map) DrawCollision(screen *ebitenv2.Image, cam *commoncamera.Camera) {
    if true { // 绘制调试边界
        p.drawBorder(screen, cam)
        p.drawPortal(screen, cam)
        p.drawBlocked(screen, cam)
        p.drawSpawnPointDebug(screen, cam) // 新增
    }
}
```

### 6. 颜色定义
- 红色 (SpawnRadius): `color.RGBA{255, 0, 0, 200}`
- 橙色 (PatrolRadius): `color.RGBA{255, 165, 0, 200}`
- 黄色 (ChaseRadius): `color.RGBA{255, 255, 0, 200}`

### 7. 注意事项
- 等距视角下，圆在屏幕上应该显示为椭圆（横向拉伸或纵向压缩），但为了简化，可以先绘制正圆，后续如有需要再调整。
- 如果要精确显示等距椭圆，需要将 radius 在 X 和 Y 方向分别缩放。

## 性能考虑
- 仅在调试模式下绘制，不影响正式版性能。
- 使用固定的 dashCount（如 48），避免过多计算。
