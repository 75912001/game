package user

import (
	"image/color"
	"saClient/src/cfg"
	"saClient/src/common"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TransitionState 过渡状态
type TransitionState int

const (
	TransitionNone    TransitionState = iota // 无过渡
	TransitionClosing                        // 关闭中 (黑幕合拢)
	TransitionOpening                        // 打开中 (黑幕展开)
)

// SceneTransition 场景过渡效果
type SceneTransition struct {
	state    TransitionState
	progress float32              // 进度 0-1
	cfg      *cfg.SceneTransition // 过渡效果配置
}

// newSceneTransition 创建场景过渡 (使用配置)
func newSceneTransition(transCfg *cfg.SceneTransition) *SceneTransition {
	return &SceneTransition{
		state:    TransitionNone,
		progress: 0,
		cfg:      transCfg,
	}
}

// In 转场进入 (黑幕展开，显示新场景)
func (t *SceneTransition) In(transCfg *cfg.SceneTransition) {
	t.cfg = transCfg
	t.state = TransitionOpening
	t.progress = 0
}

// Out 转场退出 (黑幕合拢，隐藏当前场景)
func (t *SceneTransition) Out(transCfg *cfg.SceneTransition) {
	t.cfg = transCfg
	t.state = TransitionClosing
	t.progress = 0
}

// IsActive 是否正在过渡中
func (t *SceneTransition) IsActive() bool {
	return t.state != TransitionNone
}

// Update 更新过渡状态
// 返回值:
//   - isActive: 是否正在过渡中
//   - isCloseComplete: 转场退出完成 (黑幕合拢完成)
//   - isOpenComplete: 转场进入完成 (黑幕展开完成，整个过渡结束)
func (t *SceneTransition) Update() (isActive bool, isCloseComplete bool, isOpenComplete bool) {
	if t.state == TransitionNone {
		return false, false, false
	}

	t.progress += t.cfg.Speed

	if t.progress >= 1 {
		t.progress = 1
		currentState := t.state
		t.state = TransitionNone
		t.progress = 0
		return false, currentState == TransitionClosing, currentState == TransitionOpening
	}

	return true, false, false
}

// Draw 绘制过渡效果
func (t *SceneTransition) Draw(screen *ebitenv2.Image) {
	if t.state == TransitionNone {
		return
	}

	switch t.cfg.ID {
	case cfg.SceneTransitionIDVertical:
		t.drawVertical(screen)
	case cfg.SceneTransitionIDHorizontal:
		t.drawHorizontal(screen)
	case cfg.SceneTransitionIDFade:
		t.drawFade(screen)
	default:
		t.drawVertical(screen)
	}
}

// drawVertical 绘制垂直过渡效果 (上下合拢/展开)
func (t *SceneTransition) drawVertical(screen *ebitenv2.Image) {
	screenW := float32(cfg.GCommon.ScreenMaxWidth)
	screenH := float32(cfg.GCommon.ScreenMaxHeight)
	halfH := screenH / 2

	var topHeight, bottomY float32

	if t.state == TransitionClosing { // 黑幕从上下向中间合拢
		topHeight = halfH * t.progress
		bottomY = screenH - halfH*t.progress
	} else { // 黑幕从中间向上下展开
		topHeight = halfH * (1 - t.progress)
		bottomY = screenH - halfH*(1-t.progress)
	}

	// 绘制上半部分黑幕
	if topHeight > 0 {
		vector.FillRect(screen, 0, 0, screenW, topHeight, common.Colors_Black, false)
	}

	// 绘制下半部分黑幕
	bottomHeight := screenH - bottomY
	if bottomHeight > 0 {
		vector.FillRect(screen, 0, bottomY, screenW, bottomHeight, common.Colors_Black, false)
	}
}

// drawHorizontal 绘制水平过渡效果 (左右合拢/展开)
func (t *SceneTransition) drawHorizontal(screen *ebitenv2.Image) {
	screenW := float32(cfg.GCommon.ScreenMaxWidth)
	screenH := float32(cfg.GCommon.ScreenMaxHeight)
	halfW := screenW / 2

	var leftWidth, rightX float32

	if t.state == TransitionClosing { // 黑幕从左右向中间合拢
		leftWidth = halfW * t.progress
		rightX = screenW - halfW*t.progress
	} else { // 黑幕从中间向左右展开
		leftWidth = halfW * (1 - t.progress)
		rightX = screenW - halfW*(1-t.progress)
	}

	// 绘制左半部分黑幕
	if leftWidth > 0 {
		vector.FillRect(screen, 0, 0, leftWidth, screenH, common.Colors_Black, false)
	}

	// 绘制右半部分黑幕
	rightWidth := screenW - rightX
	if rightWidth > 0 {
		vector.FillRect(screen, rightX, 0, rightWidth, screenH, common.Colors_Black, false)
	}
}

// drawFade 绘制淡入淡出效果
func (t *SceneTransition) drawFade(screen *ebitenv2.Image) {
	screenW := float32(cfg.GCommon.ScreenMaxWidth)
	screenH := float32(cfg.GCommon.ScreenMaxHeight)

	var alpha uint8

	if t.state == TransitionClosing { // 淡出 (透明度增加)
		alpha = uint8(255 * t.progress)
	} else { // 淡入 (透明度减少)
		alpha = uint8(255 * (1 - t.progress))
	}

	black := color.RGBA{R: 0, G: 0, B: 0, A: alpha}
	vector.FillRect(screen, 0, 0, screenW, screenH, black, false)
}
