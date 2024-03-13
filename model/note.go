package model

import "gorm.io/gorm"

// 存储ES, 方便搜索
type StockNote struct {
	NoteId          int
	StockId         string
	UserId          string // 用户名
	Title           string // 逻辑点
	Content         string //
	Status          int    // 是否过期(不再适用)
	Mode            int    // 1: 覆盖更新模式 2: 版本号模式(历史笔记)
	CreateTimeStamp int64
	RemindTimeStamp int64 // 提示日期
}

// mysql
type StockTag struct {
	*gorm.Model
	NoteId  int // 所属笔记
	TagName string
}
