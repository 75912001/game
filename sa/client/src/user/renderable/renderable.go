package renderable

import (
	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"sort"
)

// 所有能被渲染的物体都实现它
type Renderable interface {
	GetY() float64               // 获取用于排序的 Y 坐标 (脚底位置)
	Draw(screen *ebitenv2.Image) // 绘制方法
}

// 每一帧渲染前，对所有可见物体进行排序
func SortYAndDraw(screen *ebitenv2.Image, objects []Renderable) {
	// 根据 Y 坐标从小到大排序
	// 先画y小的
	// 后画y大的
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].GetY() < objects[j].GetY()
	})

	// 3. 按顺序绘制
	for _, obj := range objects {
		obj.Draw(screen)
	}
}
