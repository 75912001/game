package elemental

// Error 错误类型
type Error struct {
	Message string
}

func (p *Error) Error() string {
	return p.Message
}
