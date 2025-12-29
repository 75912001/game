package common

type AnimationFrame struct {
	FrameIdx  uint32 // 当前帧索引
	FrameTick uint32 // 帧计数器（用于控制动画速度）
}

func (p *AnimationFrame) Reset() {
	p.FrameIdx = 0
	p.FrameTick = 0
}

func (p *AnimationFrame) Update() {
	p.FrameTick++
	if FrameTickPerChange <= p.FrameTick {
		p.FrameTick = 0
		p.FrameIdx++
	}
}
