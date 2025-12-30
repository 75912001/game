package user

import (
	commoncamera "saClient/src/common/camera"
)

func (p *Role) Update() {
	// 更新当前场景的过渡动画
	isActive, isCloseComplete, _ := p.scene.transition.Update()

	// 转场退出完成，切换到新场景
	if isCloseComplete && p.pendingScene != nil {
		p.scene = p.pendingScene
		p.camera = commoncamera.NewCamera()
		p.sprite.wx, p.sprite.wy = p.pendingWX, p.pendingWY
		p.pendingScene = nil
		// 新场景转场进入
		p.scene.transition.In(p.scene._map.mapCfg.SceneTransition)
		p.UpdateWithAction()
		return
	}

	// 过渡动画进行中，不处理输入
	if isActive {
		p.UpdateWithAction()
		return
	}

	// 处理键盘输入
	p.HandleInput()
	p.scene.Update()
}
