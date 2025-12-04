package game

// State 游戏状态
type State int

const (
	State_StartMenu State = iota // 开始菜单
	State_Scene                  // 场景
	State_Battling               // 战斗
	State_GameOver               // 游戏结束
)
