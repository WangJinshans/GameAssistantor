package model

type GameGroupReservation struct {
	UserId         string
	GameId         string
	GameName       string // 游戏名称
	DungeonId      string // 副本id
	DungeonName    string // 副本名称
	Level          int    // 难度等级
	StartTimestamp int64  // 开始时间
	EndTimestamp   int64  // 结束时间
	IsElastic      bool   // 是否弹性
}

type GameGroupInfo struct {
	GroupId        string // 团本Id
	GroupType      string // 一拖, 二拖, 工资AA
	CaptainId      string // 团长
	GameId         string // 游戏id

	GameName       string // 游戏名称
	DungeonId      string // 副本id
	DungeonName    string // 副本名称
	Level          int    // 难度等级
	StartTimestamp int64  // 开始时间
	EndTimestamp   int64  // 结束时间
	Members        []User // 队员
}

// 认证材料, 录像回放等等
type GroupGameExtraInfo struct {
	GroupId        string // 团本Id
	UserId         string // 用户Id
	GameId         string // 游戏id
	GameName       string // 游戏名称
	DungeonId      string // 副本id
	DungeonName    string // 副本名称
	Level          int    // 难度等级
	StartTimestamp int64  // 开始时间
	EndTimestamp   int64  // 结束时间
	Message        string // 材料
}
