package camera

// Camera 摄像机
// 所有坐标均为 World 坐标系(像素)
type Camera struct {
	FollowWX int // 摄像机跟随点 (World 坐标) x
	FollowWY int // 摄像机跟随点 (World 坐标) y

	ViewportWX int // 视口左上角 (World 坐标) x
	ViewportWY int // 视口左上角 (World 坐标) y
}

func NewCamera() *Camera {
	return &Camera{}
}
