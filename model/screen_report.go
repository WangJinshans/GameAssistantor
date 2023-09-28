package model

import (
	"gorm.io/gorm"
)

type ScreenReportInfo struct {
	gorm.Model
	DeviceId  string `json: "device_id"`
	FileId    string `json:"file_id" gorm:"size:50"`
	FilePath  string `json:"file_path"`
	TimeStamp int64  `json:"timestamp"`
	Status    int    `json:"status"` // 状态 1: 未处理 2: 已处理
}
