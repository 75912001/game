package role

import (
	"saClient/src/common"
	"saClient/src/proto"
)

func (p *Role) Update() {
	// 处理键盘输入
	p.HandleInput()

	// 更新角色底部中心点坐标
	p.roleSprite.bottomCenterX = p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_BottomCenterX)
	p.roleSprite.bottomCenterY = p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_BottomCenterY)
	// 更新角色方向
	p.roleSprite.direction = p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_Direction)

	images := p.cfgRole.ResRole.Move.Frames[uint32(p.roleSprite.direction)]
	p.roleSprite.image = images[p.frameIdx%uint32(len(images))]

	frames := p.cfgRole.ResRole.Move.FrameInfo[p.roleSprite.direction]
	p.roleSprite.roleImageSprite = frames[p.frameIdx%uint32(len(frames))]

	// 更新角色中心坐标
	p.roleSprite.centerX = p.roleSprite.bottomCenterX
	p.roleSprite.centerY = p.roleSprite.bottomCenterY - p.roleSprite.roleImageSprite.Frame.Height/2

	// 更新摄像机坐标
	p.camera.X = p.roleSprite.centerX
	p.camera.Y = p.roleSprite.centerY

	// 计算摄像机屏幕坐标(左上角)
	screenX := p.camera.X - common.ScreenWidth/2
	screenY := p.camera.Y - common.ScreenHeight/2
	// 获取地图尺寸
	mapWidth, mapHeight := p.scene.GetMapSize()

	// 限制摄像机不显示地图边界之外
	if mapWidth <= common.ScreenWidth { // 地图宽度小于屏幕，地图居中
		screenX = -(common.ScreenWidth - mapWidth) / 2
	} else { // 地图宽度大于屏幕，限制边界
		if screenX < 0 {
			screenX = 0
		}
		if screenX > mapWidth-common.ScreenWidth {
			screenX = mapWidth - common.ScreenWidth
		}
	}
	if mapHeight <= common.ScreenHeight { // 地图高度小于屏幕，地图居中
		screenY = -(common.ScreenHeight - mapHeight) / 2
	} else {
		// 地图高度大于屏幕，限制边界
		if screenY < 0 {
			screenY = 0
		}
		if screenY > mapHeight-common.ScreenHeight {
			screenY = mapHeight - common.ScreenHeight
		}
	}
	p.camera.ScreenX = screenX
	p.camera.ScreenY = screenY

	p.scene.Update()
}
