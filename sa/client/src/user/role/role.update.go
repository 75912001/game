package role

import (
	"saClient/src/common"
	"saClient/src/proto"
)

func (p *Role) Update() {
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
	// 更新摄像机屏幕坐标(左上角)
	p.camera.ScreenX = p.camera.X - common.ScreenWidth/2
	p.camera.ScreenY = p.camera.Y - common.ScreenHeight/2

	p.scene.Update()
}
