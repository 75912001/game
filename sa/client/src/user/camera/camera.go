package camera

type Camera struct {
	X int // 摄像机位置X
	Y int // 摄像机位置Y
}

func NewCamera(x, y int) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}
