package user

func (p *Role) Update() {
	// 处理键盘输入
	p.HandleInput()
	p.scene.Update()
}
