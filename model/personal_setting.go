package model

import "gorm.io/gorm"

type PersonalSettingInfo struct {
	gorm.Model
	UserId string // 用户id
	Key    string // 设置项(note_mode: 笔记更新模式, 1: 覆盖更新, 2: 多版本模式)
	Value  string // 值
}

func (s *PersonalSettingInfo) TableName() string {
	return "personal_setting_info"
}
