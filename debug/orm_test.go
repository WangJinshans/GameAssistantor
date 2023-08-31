package debug

import (
	"game_assistantor/model"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	"testing"
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"index" json:"name"`
	UuID      uint
	// 使用默认外键在写入User的时候默认会把User表中的ID字段写入到Company表中的UserID字段
	// 可使用foreignKey,references指定关联关系, foreignKey:UserName;references:name 会将User表中的name字段写到Company表中的UserName字段
	// foreignKey是其他表字段, references是本表字段

	//CompanyReferID uint
	Company    Company    `gorm:"foreignKey:CompanyReferID"`       // belong to 适用于关联id不在关联模型的表中, 关联模型是一个独立的模型
	CreditCard CreditCard `gorm:"foreignKey:UuID;references:UuID"` // has one, 映射关系会被写到credit_card表中
}

type Company struct {
	gorm.Model
	CompanyReferID uint
	Name           string `json:"name"`
	Desc           string `json:"desc"`
}

type CreditCard struct {
	gorm.Model
	Number string
	//UserID uint
	UuID uint
}

func TestGormFilePartition(t *testing.T) {

	engine, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Info().Msgf("create db error: %v", err)
		return
	}
	engine.AutoMigrate(&model.FilePartition{}, &model.PartitionInfo{})

	p1 := model.PartitionInfo{
		SegmentId:   "segment1",
		SegmentName: "segment1",
		SegmentPath: "zzzzz",
	}

	p2 := model.PartitionInfo{
		SegmentId:   "segment2",
		SegmentName: "segment2",
		SegmentPath: "xxxx",
	}
	f := model.FilePartition{
		FileId:   "1111",
		FileName: "1111.mp4",
		FilePath: "xxxxx/xxxx",
		PartitionList: []model.PartitionInfo{
			p1,
			p2,
		},
	}
	engine.Save(&f)
}

func TestGorm(t *testing.T) {

	engine, err := gorm.Open(sqlite.Open("./gorm.db"), &gorm.Config{})
	if err != nil {
		log.Info().Msgf("create db error: %v", err)
		return
	}
	engine.AutoMigrate(&Company{}, &User{}, &CreditCard{})
	//company := Company{
	//	Name: "Company1",
	//}
	//
	//creditCard := CreditCard{
	//	Number: "CreditCard1",
	//}
	//
	//user := User{
	//	Name:           "XXXX",
	//	Company:        company,
	//	CreditCard:     creditCard,
	//}
	//
	//engine.Save(&user)

	//var u User
	//err = engine.Table("companies").Select("users.*").Joins("left join users on companies.id = users.cid").Where("companies.id = ?", 1).Find(&u).Error
	//log.Info().Msgf("user is: %v, error is: %v", u, err)
}
