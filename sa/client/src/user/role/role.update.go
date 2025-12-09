package role

import "saClient/src/proto"

func (p *Role) Update() {
	p.roleSprite.x = p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_X)
	p.roleSprite.y = p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_Y)
	p.roleSprite.direction = p.GetValueInt(proto.AssetIDRecord_AssetIDRecord_Direction)

	images := p.cfgRole.ResRole.Move.Frames[uint32(p.roleSprite.direction)]
	p.roleSprite.image = images[p.frameIdx%uint32(len(images))]

	frames := p.cfgRole.ResRole.Move.FrameInfo[p.roleSprite.direction]
	p.roleSprite.roleImageSprite = frames[p.frameIdx%uint32(len(frames))]

	p.camera.X = p.roleSprite.x + p.roleSprite.roleImageSprite.Frame.Width/2
	p.camera.Y = p.roleSprite.y + p.roleSprite.roleImageSprite.Frame.Height/2

	p.scene.Update()
}
