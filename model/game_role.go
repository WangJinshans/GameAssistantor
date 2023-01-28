package model

type GameAccount struct {
	AccountType   string
	AccountId     string
	AccountPwd    string
	AccountStatus string // 账号状态
}

type DNFRole struct {
	AccountId    string // 账号
	RoleId       string
	RoleName     string
	SummonerType string // 游戏内角色类型
	Level        int    // 等级
	Reputation   string // 名望
	Backup       string // 备注
	Money        int64  // 金币
}
