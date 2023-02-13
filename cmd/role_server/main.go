package main

import (
	"fmt"
	"game_assistantor/api/login"
	"game_assistantor/api/v1/game_account"
	"game_assistantor/config"
	"game_assistantor/middlerware"
	"game_assistantor/repository"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	conf        config.AssistantConfig
	engine      *gorm.DB
	redisClient *redis.Client
)

func init() {
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Error().Msgf("fail to decode config.toml, error is: %v", err)
		return
	}
	initDatabase()
}

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
}

func SyncTables() (err error) {
	err = engine.AutoMigrate(
	)
	return
}

func StartServer() {
	r := gin.New()
	r.Use(middlerware.Cors())
	r.Use(gin.Recovery())
	r.Static("/static", "../../static/")
	r.GET("/get_public_key", login.GetPublicKey)
	r.POST("/login", login.Login)
	r.POST("/logout", login.Logout)
	r.POST("/register", login.Register)
	r.POST("/refresh_token", login.RefreshToken)
	r.GET("/get_qrcode", login.InitQrCode)
	r.GET("/get_qrcode_status", login.QueryQrCode)
	r.POST("/set_qrcode_status", login.SetQrCodeStatus)
	r.POST("/scan_qrcode", login.ScanQrCode)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("v1")
	{
		ecoGroup := v1.Group("role")
		{
			ecoGroup.POST("/add", game_account.GameRoleApi.GetAccountInfo)
		}
	}
	r.Run(":8088")
}

func main() {
	//StartServer()
}



