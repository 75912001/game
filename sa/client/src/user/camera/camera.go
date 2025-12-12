package camera

// Camera 摄像机
// 所有坐标均为 World 坐标系(像素)
type Camera struct {
	// 摄像机跟随点 (World 坐标)
	FollowX int
	FollowY int

	// 视口左上角 (World 坐标)
	ViewportX int
	ViewportY int
}

func NewCamera() *Camera {
	return &Camera{}
}
