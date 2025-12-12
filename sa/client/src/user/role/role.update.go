package role

import (
	"saClient/src/cfg"
	"saClient/src/proto"
)

func (p *Role) Update() {
	// 处理键盘输入
	p.HandleInput()

	// 初始化 World 坐标 (首次或重生时从 Tile 坐标转换)
	if p.roleSprite.bottomCenterWorldX == 0 && p.roleSprite.bottomCenterWorldY == 0 {
		tx := p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TX)
		ty := p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_BottomCenter_TY)
		p.roleSprite.bottomCenterWorldX, p.roleSprite.bottomCenterWorldY = p.scene.TileToWorld(float64(tx), float64(ty))
	}

	// 更新角色方向
	p.roleSprite.orientation = p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_Orientation)

	images := p.cfgRole.ResRole.Move.Frames[uint32(p.roleSprite.orientation)]
	p.roleSprite.image = images[p.frameIdx%uint32(len(images))]

	frames := p.cfgRole.ResRole.Move.FrameInfo[p.roleSprite.orientation]
	p.roleSprite.roleImageSprite = frames[p.frameIdx%uint32(len(frames))]

	// 更新角色中心 World 坐标 (脚底中心向上偏移半个角色高度)
	p.roleSprite.centerWorldX = p.roleSprite.bottomCenterWorldX
	p.roleSprite.centerWorldY = p.roleSprite.bottomCenterWorldY - float64(p.roleSprite.roleImageSprite.Frame.Height/2)

	// 更新摄像机跟随点 (World 坐标)
	p.camera.FollowX = int(p.roleSprite.centerWorldX)
	p.camera.FollowY = int(p.roleSprite.centerWorldY)

	// 计算摄像机视口左上角 (World 坐标)
	viewportX := p.camera.FollowX - cfg.GCommon.ScreenMaxWidth/2
	viewportY := p.camera.FollowY - cfg.GCommon.ScreenMaxHeight/2

	// 获取地图尺寸
	mapWidth, mapHeight := p.scene.GetMapPixeSize()

	// 限制摄像机不显示地图边界之外
	if mapWidth <= cfg.GCommon.ScreenMaxWidth {
		viewportX = -(cfg.GCommon.ScreenMaxWidth - mapWidth) / 2
	} else {
		if viewportX < 0 {
			viewportX = 0
		}
		if viewportX > mapWidth-cfg.GCommon.ScreenMaxWidth {
			viewportX = mapWidth - cfg.GCommon.ScreenMaxWidth
		}
	}

	if mapHeight <= cfg.GCommon.ScreenMaxHeight {
		viewportY = -(cfg.GCommon.ScreenMaxHeight - mapHeight) / 2
	} else {
		if viewportY < 0 {
			viewportY = 0
		}
		if viewportY > mapHeight-cfg.GCommon.ScreenMaxHeight {
			viewportY = mapHeight - cfg.GCommon.ScreenMaxHeight
		}
	}

	p.camera.ViewportX = viewportX
	p.camera.ViewportY = viewportY

	p.scene.Update()
}
