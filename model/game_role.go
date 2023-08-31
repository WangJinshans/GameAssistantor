package model

import "gorm.io/gorm"

type GameAccount struct {
	gorm.Model
	AccountType   string // 账号类型
	AccountId     string
	AccountPwd    string
	AccountStatus string // 账号状态
}

type GameRole struct {
	AccountId    string // 账号
	RoleId       string
	RoleName     string
	SummonerType string // 游戏内角色类型
}

type DNFRole struct {
	GameRole
	Level      int    // 等级
	Reputation string // 名望
	Backup     string // 备注
	Money      int64  // 金币
}
