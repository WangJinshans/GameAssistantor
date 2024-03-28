package main

import (
	"fmt"
	"game_assistantor/api/login"
	"game_assistantor/auth"
	"game_assistantor/config"
	"game_assistantor/middlerware"
	"game_assistantor/model"
	"game_assistantor/repository"
	"game_assistantor/route"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	conf   config.AssistantConfig
	engine *gorm.DB
)

func init() {
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	initDatabase()
}

func initDatabase() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Mysql.UserName, conf.Mysql.Password, conf.Mysql.Address, conf.Mysql.Database)

	log.Info().Msgf("connection string is: %s", dsn)
	engine, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = engine.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	repository.SetupEngine(engine)
}

func StartServer() {
	r := gin.Default()
	r.Use(middlerware.Cors())
	r.Use(gin.Recovery())

	api := r.Group(route.RootTPath)
	base := api.Group(route.BaseGroupName)
	{
		base.GET(route.PublicKeyPath, login.GetPublicKey)
		base.POST(route.LoginPath, login.Login)
		base.POST(route.LogoutPath, login.Logout)
		base.POST(route.RegisterPath, login.Register)
		base.POST(route.TokenPath, login.RefreshToken)
		base.POST(route.QrCodePath, login.InitQrCode)
		base.GET(route.QrCodeStatusPath, login.QueryQrCode)
		base.PATCH(route.QrCodeStatusPath, login.SetQrCodeStatus)
		base.POST(route.QrCodeScanPath, login.ScanQrCode)
	}

	authGroup := api.Group(route.RouterAuthGroupName)
	{
		authGroup.POST(route.VerifyPath, auth.VerifyAccess)
	}

	r.Run(":8088")
}

func main() {
	StartServer()
}
