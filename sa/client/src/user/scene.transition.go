package user

import (
	"image/color"
	"saClient/src/cfg"

	ebitenv2 "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TransitionState 过渡状态
type TransitionState int

const (
	TransitionNone    TransitionState = iota // 无过渡
	TransitionClosing                        // 关闭中 (黑幕从上下向中间合拢)
	TransitionOpening                        // 打开中 (黑幕从中间向上下展开)
)

// SceneTransition 场景过渡效果
type SceneTransition struct {
	state    TransitionState
	progress float32 // 进度 0-1
	speed    float32 // 每帧增加的进度
}

// newSceneTransition 创建场景过渡
func newSceneTransition() *SceneTransition {
	return &SceneTransition{
		state:    TransitionNone,
		progress: 0,
		speed:    0.1,
	}
}

// TransitionIn 转场进入 (黑幕从中间向上下展开，显示新场景)
func (t *SceneTransition) TransitionIn() {
	t.state = TransitionOpening
	t.progress = 0
}

// TransitionOut 转场退出 (黑幕从上下向中间合拢，隐藏当前场景)
func (t *SceneTransition) TransitionOut() {
	t.state = TransitionClosing
	t.progress = 0
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

	t.progress += t.speed

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

	screenW := float32(cfg.GCommon.ScreenMaxWidth)
	screenH := float32(cfg.GCommon.ScreenMaxHeight)
	halfH := screenH / 2

	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	var topHeight, bottomY float32

	if t.state == TransitionClosing { // 关闭中: 黑幕从上下向中间合拢
		topHeight = halfH * t.progress
		bottomY = screenH - halfH*t.progress
	} else { // 打开中: 黑幕从中间向上下展开
		topHeight = halfH * (1 - t.progress)
		bottomY = screenH - halfH*(1-t.progress)
	}

	// 绘制上半部分黑幕
	if topHeight > 0 {
		vector.FillRect(screen, 0, 0, screenW, topHeight, black, false)
	}

	// 绘制下半部分黑幕
	bottomHeight := screenH - bottomY
	if bottomHeight > 0 {
		vector.FillRect(screen, 0, bottomY, screenW, bottomHeight, black, false)
	}
}

// IsActive 是否正在过渡中
func (t *SceneTransition) IsActive() bool {
	return t.state != TransitionNone
}

// IsClosing 是否正在关闭中
func (t *SceneTransition) IsClosing() bool {
	return t.state == TransitionClosing
}

// IsOpening 是否正在打开中
func (t *SceneTransition) IsOpening() bool {
	return t.state == TransitionOpening
}
