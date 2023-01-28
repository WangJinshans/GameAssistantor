package model

// 每天根据提示生成工作任务
type TipMessage struct {
	TipId            string
	Content          string
	TipType          string // 角色提示 账号提示
	ExpiredTimeStamp int64  // 到期时间
}
