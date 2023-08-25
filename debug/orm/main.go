package main

import (
	"fmt"
	"game_assistantor/config"
	"game_assistantor/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	engine *gorm.DB
	conf   config.AssistantConfig
)

type Company struct {
	gorm.Model
	ID   int
	Name string
}

type User struct {
	// has one, belongs to 都必须存在外键, 默认规则: 类型名 + 主键
	// `gorm:"foreignkey:CompanyId;references:ID"`
	// reference 指向Company中字段, foreignkey 指向本结构字段
	gorm.Model
	Name      string  `gorm:"size:20"`
	CompanyId int     // belongs to
	Company   Company /*`gorm:"foreignkey:CompanyId;references:ID"`*/

	CreditCard []CreditCard `gorm:"foreignkey:Name;references:Number"` // has one
}

type CreditCard struct {
	gorm.Model
	Number string `gorm:"size:20;index"`
	UserID int
}

func init() {
	//_, err := toml.DecodeFile("config.toml", &conf)
	//if err != nil {
	//	log.Error().Msgf("fail to decode config.toml, error is: %v", err)
	//	return
	//}
	//initDatabase()
}

// reference: https://blog.csdn.net/weixin_46618592/article/details/127194231
func initDatabase() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Mysql.UserName, conf.Mysql.Password, conf.Mysql.Address, conf.Mysql.Database)
	engine, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = SyncTables()
	if err != nil {
		panic(err)
	}
	repository.SetupEngine(engine)
	//c := Company{
	//	ID:   1,
	//	Name: "xxxxx",
	//}
	//credit := CreditCard{
	//	Number: "Lance",
	//	UserID: 333,
	//}
	//credit1 := CreditCard{
	//	Number: "Lance",
	//	UserID: 11,
	//}
	//u := User{
	//	Name:       "Lance",
	//	Company:    c,
	//	CreditCard: []CreditCard{credit, credit1},
	//}
	////engine.Save(&c)
	//err = engine.Create(&u).Error
	//log.Info().Msgf("error is: %v", err)
	//c := Company{}
	//c.Name = "xxxxx"
	////engine.Model(&Company{}).Where("name = ?", "xxxxx").First(&c)
	//u1 := User{}
	//u1.Name = "Lance"
	//engine.Model(&User{}).Preload("Company").Preload("CreditCard").Take(&u1)
	//log.Info().Msgf("u1 is: %v", u1.CreditCard)
}

func SyncTables() (err error) {
	err = engine.AutoMigrate(
		//&model.GameAccount{},
		&Company{},
		&User{},
		&CreditCard{},
	)
	//err = engine.Migrator().DropColumn(&User{},"company_refer")
	//log.Info().Msgf("err is: %v", err)
	return
}

func main() {

}
