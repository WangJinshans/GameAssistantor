package main

import (
	"context"
	"fmt"
	"game_assistantor/api/role"
	"game_assistantor/api/v1/device"
	"game_assistantor/api/v1/note"
	"game_assistantor/api/v1/user"
	"game_assistantor/config"
	_ "game_assistantor/docs" // 引入文档
	"game_assistantor/middlerware"
	"game_assistantor/model"
	"game_assistantor/repository"
	"game_assistantor/route"
	"game_assistantor/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	engine      *gorm.DB
	signals     chan os.Signal
	conf        config.AssistantConfig
	redisClient *redis.Client
)

func init() {
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Error().Msgf("fail to decode config.toml, error is: %v", err)
		return
	}

	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	initDatabase()
	initRedis()
}

func initRedis() {
	redisClient = utils.GetRedisClient(conf.Redis.Address, conf.Redis.DB, conf.Redis.Password)
	repository.SetupRedisClient(redisClient)
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
		&model.ScreenReportInfo{},
	)
	//err = engine.Migrator().DropColumn(&GameAccount{},"xxxx_xxx")
	//log.Info().Msgf("err is: %v", err)
	return
}

func StartWebServer(ctx context.Context) {

	r := gin.New()
	r.Use(middlerware.Cors())
	r.Use(gin.Recovery())
	r.Use(middlerware.RequestInfo())

	r.Static("/static", "../../static/")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // swag init -g ./cmd/assistantor_server

	api := r.Group(route.RouterApiGroupName)
	{
		v1 := api.Group(route.RouterVersionGroupName)
		{
			// v1.Use(middlerware.AuthMiddleWare()) // 权限验证 请求服务认证
			roleGroup := v1.Group(route.RoleGroupName)
			{
				roleGroup.POST(route.RolePath, role.RoleApi.AddRoleForUser)
				roleGroup.GET(route.RoleSPath, role.RoleApi.GetAllRoles)
				roleGroup.DELETE(fmt.Sprintf("%s:account_id", route.RolePath), role.RoleApi.DeleteRole)
			}
			userGroup := v1.Group(route.UserGroupName)
			{
				userGroup.GET(route.UsersPath, user.UserApi.GetUsersInfo)
				userGroup.GET(fmt.Sprintf("%s:user_id", route.UserPath), user.UserApi.GetUserInfo)
				userGroup.PATCH(fmt.Sprintf("%s:user_id", route.UserPath), user.UserApi.UpdateUserInfo)
				userGroup.PATCH(fmt.Sprintf("%s:user_id", route.UserPasswordPath), user.UserApi.UpdateUserPassword)
			}
			noteGroup := v1.Group(route.NoteGroupName)
			{
				noteGroup.GET(fmt.Sprintf("%s:user_id/:note_id", route.NotePath), note.NoteApi.GetNote)
				noteGroup.GET(fmt.Sprintf("%s/:user_id", route.NotesPath), note.NoteApi.GetNoteList)
				noteGroup.POST(fmt.Sprintf("%s", route.NoteCreatePath), note.NoteApi.CreateNote)
				noteGroup.PATCH(fmt.Sprintf("%s:user_id/:note_id", route.NotePath), note.NoteApi.UpdateNote)
			}
			deviceGroup := v1.Group(route.DeviceGroupName)
			{
				deviceGroup.GET(route.DevicesPath, device.DeviceApi.GetDevicesList)
				deviceGroup.POST("/command", device.DeviceApi.SendDeviceCommand)
				deviceGroup.POST("/report", device.DeviceApi.StatusReport)
				deviceGroup.POST("/set_report", device.DeviceApi.SetReportStatus)
				deviceGroup.GET("/get_reports", device.DeviceApi.GetReportList)
			}
		}
	}
	r.Run(":8088")
}

func main() {

	ctx := context.Background()
	cancleCtx, cancleFunc := context.WithCancel(ctx)

	go StartWebServer(cancleCtx)
	// go service.StartDeviceService(cancleCtx) // 启动设备服务

	<-signals
	cancleFunc()
	// 剩下2s等待退出
	time.Sleep(2 * time.Second)
}
