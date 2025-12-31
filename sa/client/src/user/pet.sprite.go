package user

import (
	"saClient/src/proto"
	"saClient/src/res"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
)

// PetSprite 宠物-精灵
type PetSprite struct {
	wx float32 // 脚底中心 World 坐标 X
	wy float32 // 脚底中心 World 坐标 Y

	tx float32 // 脚底中心 Tile 坐标 X
	ty float32 // 脚底中心 Tile 坐标 Y

	centerWX float32 // 中心 World 坐标 X
	centerWY float32 // 中心 World 坐标 Y

	action      proto.PetAction        // 当前播放的动作 (Move, Attack等)
	orientation proto.AssetOrientation // 方向
	image       *ebitenv2.Image        // 图片
	imageSprite *res.PetImageSprite    // 图片配置
}

// GetWX 获取脚底中心 World 坐标 X
func (p *PetSprite) GetWX() float32 {
	return p.wx
}

// GetWY 获取脚底中心 World 坐标 Y
func (p *PetSprite) GetWY() float32 {
	return p.wy
}

// SetWX 设置脚底中心 World 坐标 X
func (p *PetSprite) SetWX(wx float32) {
	p.wx = wx
}

// SetWY 设置脚底中心 World 坐标 Y
func (p *PetSprite) SetWY(wy float32) {
	p.wy = wy
}

// GetCenterWX 获取中心 World 坐标 X
func (p *PetSprite) GetCenterWX() float32 {
	return p.centerWX
}

// GetCenterWY 获取中心 World 坐标 Y
func (p *PetSprite) GetCenterWY() float32 {
	return p.centerWY
}

// SetCenterWX 设置中心 World 坐标 X
func (p *PetSprite) SetCenterWX(centerWX float32) {
	p.centerWX = centerWX
}

// SetCenterWY 设置中心 World 坐标 Y
func (p *PetSprite) SetCenterWY(centerWY float32) {
	p.centerWY = centerWY
}

// GetAction 获取当前动作
func (p *PetSprite) GetAction() proto.PetAction {
	return p.action
}

// SetAction 设置当前动作
func (p *PetSprite) SetAction(action proto.PetAction) {
	p.action = action
}

// GetOrientation 获取方向
func (p *PetSprite) GetOrientation() proto.AssetOrientation {
	return p.orientation
}

// SetOrientation 设置方向
func (p *PetSprite) SetOrientation(orientation proto.AssetOrientation) {
	p.orientation = orientation
}
