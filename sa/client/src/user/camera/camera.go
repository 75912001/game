package camera

type Camera struct {
	X int // 摄像机位置X
	Y int // 摄像机位置Y

	ScreenX int // 屏幕左上角X
	ScreenY int // 屏幕左上角Y
}

func NewCamera() *Camera {
	return &Camera{}
}
