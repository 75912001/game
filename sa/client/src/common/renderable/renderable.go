package renderable

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	commoncamera "saClient/src/common/camera"
	"sort"
)

// 所有能被渲染的物体都实现它
type IRenderable interface {
	GetWY() float32                                           // 获取用于排序的 Y 坐标
	Draw(screen *ebitenv2.Image, camera *commoncamera.Camera) // 绘制方法
}

// 每一帧渲染前，对所有可见物体进行排序
/*
❗致命陷阱：❌千万不要在 sort.Slice 中使用 Epsilon！
这是游戏开发中常见的坑。如果你打算把这个 Float32Less 用在 sort.Slice 的排序函数里：
// ⚠️ 警告：这样做可能会导致排序崩溃或结果不确定
sort.Slice(objects, func(i, j int) bool {
    return Float32Less(objects[i].Y, objects[j].Y)
})

为什么不行？ 排序算法（如 Go 的 pdqsort）要求比较函数满足传递性。如果引入了 Epsilon（模糊比较），就会破坏“相等的传递性”。
举例： 设 epsilon = 1.0
A = 1.0, B = 1.8, C = 2.6
A 和 B 的差是 0.8（小于 1.0），被判定为相等。
B 和 C 的差是 0.8（小于 1.0），被判定为相等。
但是，A 和 C 的差是 1.6（大于 1.0），被判定为不相等（A < C）。
结果： 排序算法会懵圈：A等于B，B等于C，为什么A不等于C？这会导致排序结果在不同帧之间随机跳动（Z-Fighting），甚至某些算法会死循环。

✅正确的建议

直接使用原生操作符，不要用 Epsilon。
sort.Slice(objects, func(i, j int) bool {
    // 简单粗暴，性能最高，且符合排序算法的数学要求
    return objects[i].Y < objects[j].Y
})
理由：对于渲染遮挡，如果两个物体 Y 坐标只差 0.000001，谁在前谁在后肉眼根本看不出来，直接按计算机的微小差异排就行，不需要人为容错。
*/
func SortYAndDraw(screen *ebitenv2.Image, camera *commoncamera.Camera, objects []IRenderable) {
	// 根据 Y 坐标从小到大排序 先画y小的, 后画y大的
	sort.Slice(objects, func(i, j int) bool {
		// ❌ return xutil.Float32Equal(objects[i].GetWY(), objects[j].GetWY())
		return objects[i].GetWY() < objects[j].GetWY()
	})
	// 按顺序绘制
	for _, obj := range objects {
		obj.Draw(screen, camera)
	}
}
