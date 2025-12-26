package user

import (
	"math"
	"saClient/src/cfg"

	astar "github.com/beefsack/go-astar"
)

// AStarWorldPoint World 坐标点
type AStarWorldPoint struct {
	WX, WY float32
}

// astarNode A* 搜索节点 (实现 astar.Pather 接口)
type astarNode struct {
	GX, GY int
	Grid   *cfg.TiledMapLogicalGridMgr
	World  *astarNodeWorld
}

// astarNodeWorld 节点世界
type astarNodeWorld struct {
	Grid         *cfg.TiledMapLogicalGridMgr
	Nodes        map[int]*astarNode
	SearchCount  int // 搜索计数
	MaxSearchNum int // 最大搜索数
}

func newAstarNodeWorld(grid *cfg.TiledMapLogicalGridMgr, maxSearch int) *astarNodeWorld {
	return &astarNodeWorld{
		Grid:         grid,
		Nodes:        make(map[int]*astarNode),
		MaxSearchNum: maxSearch,
	}
}

func (w *astarNodeWorld) getNode(gx, gy int) *astarNode {
	key := gy*w.Grid.GridW + gx
	if node, ok := w.Nodes[key]; ok {
		return node
	}
	node := &astarNode{GX: gx, GY: gy, Grid: w.Grid, World: w}
	w.Nodes[key] = node
	return node
}

// 8 方向
var astarDirs = [8]struct {
	dx, dy int
	cost   float64
}{
	{0, -1, 1.0},
	{0, 1, 1.0},
	{-1, 0, 1.0},
	{1, 0, 1.0},
	{-1, -1, math.Sqrt2},
	{1, -1, math.Sqrt2},
	{-1, 1, math.Sqrt2},
	{1, 1, math.Sqrt2},
}

// PathNeighbors 实现 astar.Pather
func (n *astarNode) PathNeighbors() []astar.Pather {
	// 超过搜索限制，返回空邻居强制结束
	n.World.SearchCount++
	if n.World.SearchCount > n.World.MaxSearchNum {
		return nil
	}

	var neighbors []astar.Pather
	for _, dir := range astarDirs {
		nx, ny := n.GX+dir.dx, n.GY+dir.dy

		if n.Grid.IsBlocked(nx, ny) {
			continue
		}

		// 对角线防穿墙
		if dir.dx != 0 && dir.dy != 0 {
			if n.Grid.IsBlocked(n.GX+dir.dx, n.GY) && n.Grid.IsBlocked(n.GX, n.GY+dir.dy) {
				continue
			}
		}

		neighbors = append(neighbors, n.World.getNode(nx, ny))
	}
	return neighbors
}

// PathNeighborCost 实现 astar.Pather
func (n *astarNode) PathNeighborCost(to astar.Pather) float64 {
	t := to.(*astarNode)
	if n.GX != t.GX && n.GY != t.GY {
		return math.Sqrt2
	}
	return 1.0
}

// PathEstimatedCost 实现 astar.Pather
func (n *astarNode) PathEstimatedCost(to astar.Pather) float64 {
	t := to.(*astarNode)
	dx := math.Abs(float64(t.GX - n.GX))
	dy := math.Abs(float64(t.GY - n.GY))
	if dx > dy {
		return dx + (math.Sqrt2-1)*dy
	}
	return dy + (math.Sqrt2-1)*dx
}

// AStarPathfinder A* 寻路器
type AStarPathfinder struct {
	grid         *cfg.TiledMapLogicalGridMgr
	maxSearchNum int
}

// NewAStarPathfinder 创建寻路器
func NewAStarPathfinder(grid *cfg.TiledMapLogicalGridMgr) *AStarPathfinder {
	return &AStarPathfinder{
		grid:         grid,
		maxSearchNum: 5000000, // 增大搜索限制
	}
}

// SetMaxSearchNum 设置最大搜索数
func (p *AStarPathfinder) SetMaxSearchNum(num int) {
	p.maxSearchNum = num
}

// FindPath 寻找路径
func (p *AStarPathfinder) FindPath(startWX, startWY, endWX, endWY float32) []AStarWorldPoint {
	if p.grid == nil {
		return nil
	}

	startGX, startGY := p.grid.W2Grid(startWX, startWY)
	endGX, endGY := p.grid.W2Grid(endWX, endWY)

	if p.grid.IsBlocked(endGX, endGY) {
		return nil
	}

	if startGX == endGX && startGY == endGY {
		return []AStarWorldPoint{{endWX, endWY}}
	}

	world := newAstarNodeWorld(p.grid, p.maxSearchNum)
	startNode := world.getNode(startGX, startGY)
	endNode := world.getNode(endGX, endGY)

	path, _, found := astar.Path(startNode, endNode)
	if !found || len(path) == 0 {
		return nil
	}

	// 路径转换 (go-astar 返回终点→起点，需反转，不含起点)
	worldPath := make([]AStarWorldPoint, 0, len(path)-1)
	for i := len(path) - 2; i >= 0; i-- {
		node := path[i].(*astarNode)
		grid := p.grid.GetGrid(node.GX, node.GY)
		worldPath = append(worldPath, AStarWorldPoint{grid.WX, grid.WY})
	}

	if len(worldPath) > 0 {
		worldPath[len(worldPath)-1] = AStarWorldPoint{endWX, endWY}
	}

	return worldPath
}
