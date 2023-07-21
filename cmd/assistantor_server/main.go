package main

import (
	"fmt"
	"game_assistantor/api/login"
	"game_assistantor/api/role"
	"game_assistantor/api/v1/game_account"
	"game_assistantor/api/v1/user"
	"game_assistantor/config"
	_ "game_assistantor/docs" // 引入文档
	"game_assistantor/middlerware"
	"game_assistantor/model"
	"game_assistantor/repository"
	"game_assistantor/route"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		&model.GameAccount{},
	)
	//err = engine.Migrator().DropColumn(&GameAccount{},"xxxx_xxx")
	//log.Info().Msgf("err is: %v", err)
	return
}

func StartServer() {

	r := gin.New()
	r.Use(middlerware.Cors())
	r.Use(gin.Recovery())
	r.Use(middlerware.RequestInfo())

	r.Static("/static", "../../static/")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // swag init -g ./cmd/assistantor_server

	v1 := r.Group(route.RouterVersionGroupName)
	{
		base := v1.Group(route.BaseGroupName)
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
		roleGroup := v1.Group(route.RoleGroupName)
		{
			roleGroup.POST(route.RolePath, role.RoleApi.AddRoleForUser)
			roleGroup.GET(route.RoleSPath, role.RoleApi.GetAllRoles)
			roleGroup.DELETE(fmt.Sprintf("%s:account_id", route.RolePath), role.RoleApi.DeleteRole)
		}
		gameRoleGroup := v1.Group(route.GameRoleGroupName)
		{
			gameRoleGroup.POST(route.GameRolePath, game_account.GameRoleApi.AddAccount)
			gameRoleGroup.GET(fmt.Sprintf("%s:account_id", route.GameRolePath), game_account.GameRoleApi.GetAccountInfo)
			gameRoleGroup.PATCH(fmt.Sprintf("%s:account_id", route.GameRolePath), game_account.GameRoleApi.UpdateAccountInfo)
		}
		userGroup := v1.Group(route.UserGroupName)
		{
			userGroup.GET(route.UsersPath, user.UserApi.GetUsersInfo)
			userGroup.GET(fmt.Sprintf("%s:user_id", route.UserPath), user.UserApi.GetUserInfo)
			userGroup.PATCH(fmt.Sprintf("%s:user_id", route.UserPath), user.UserApi.UpdateUserInfo)
			userGroup.PATCH(fmt.Sprintf("%s:user_id", route.UserPasswordPath), user.UserApi.UpdateUserPassword)
		}
	}
	r.Run(":8088")

}

func main() {
	StartServer()
}
