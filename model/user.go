package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserId      string `json:"user_id" gorm:"primarykey;index;size:32"`
	UserName    string `json:"user_name"`
	NickName    string `json:"nick_name"`
	PassWord    string `json:"pass_word"`
	EmailAddr   string `json:"email_addr"`
	PhoneNumber string `json:"phone_number"`
	IdCard      string `json:"id_card"` // 身份信息
	LevelExpire int64  `json:"level_expire"`
	AddressID   uint   `json:"address_id"` // 地址
}
