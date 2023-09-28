package device

import (
	"errors"
	"fmt"
	"game_assistantor/common"
	"game_assistantor/model"
	"game_assistantor/repository"
	"game_assistantor/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var DeviceApi ApiDevice

type ApiDevice struct {
}

func (*ApiDevice) GetDevicesList(ctx *gin.Context) {

	log.Info().Msgf("get device list...")
	deviceList, err := service.GetDeviceList()
	if err != nil {
		log.Info().Msgf("fail to get user, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": deviceList,
	})
}

// 直线移动 web端
func (d *ApiDevice) SendDeviceCommand(ctx *gin.Context) {

	type parameter struct {
		DeviceId  string `json:"device_id"`
		Direction string `json:"direction"`
		Duration  int    `json:"duration"`
		Reset     string `json:"reset"`
	}

	var param parameter

	err := ctx.BindJSON(&param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}

	log.Info().Msgf("param is: %#v", param)
	if param.DeviceId == "" {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "device id error",
		})
		return
	}

	if param.Reset == "reset" {
		command, err := d.generateCommand(param.DeviceId, "release_all", "reset")
		if err != nil {
			ctx.JSON(200, gin.H{
				"code":    common.Fail,
				"message": err.Error(),
			})
			return
		}

		err = service.SendCommand(param.DeviceId, command)
		if err != nil {
			log.Error().Msgf("device id: %s, fail to send command, error is: %v", param.DeviceId, err)
			ctx.JSON(200, gin.H{
				"code":    common.Fail,
				"message": err.Error(),
			})
		}
	} else {
		command, err := d.generateCommand(param.DeviceId, param.Direction, "down")
		if err != nil {
			ctx.JSON(200, gin.H{
				"code":    common.Fail,
				"message": err.Error(),
			})
			return
		}

		err = service.SendCommand(param.DeviceId, command)
		if err != nil {
			log.Error().Msgf("device id: %s, fail to send command, error is: %v", param.DeviceId, err)
			ctx.JSON(200, gin.H{
				"code":    common.Fail,
				"message": err.Error(),
			})
		}

		time.Sleep(time.Millisecond * time.Duration(param.Duration*100)) // 1即为0.1s, 与脚本端保持一致

		command, err = d.generateCommand(param.DeviceId, param.Direction, "up")
		if err != nil {
			ctx.JSON(200, gin.H{
				"code":    common.Fail,
				"message": err.Error(),
			})
			return
		}

		err = service.SendCommand(param.DeviceId, command)
		if err != nil {
			log.Error().Msgf("device id: %s, fail to send command, error is: %v", param.DeviceId, err)
			ctx.JSON(200, gin.H{
				"code":    common.Fail,
				"message": err.Error(),
			})
		}
	}

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
}

func (*ApiDevice) generateCommand(deviceId string, direction string, operation string) (command []byte, err error) {
	// command = []byte(fmt.Sprintf("#%s#reset_all#", deviceId))
	transformedKey, ok := common.KeyBoardMap[direction]
	if !ok {
		err = errors.New("unsupport key")
		return
	}
	command = []byte(fmt.Sprintf("#device_server@publisher@%s@%s@%s@end_string#", deviceId, transformedKey, operation))
	return
}

func (*ApiDevice) StatusReport(ctx *gin.Context) {

	param, err := ctx.MultipartForm()
	if err != nil {
		log.Info().Msgf("fail to decode params, error is: %v", err)
		return
	}
	deviceId := param.Value["device_id"][0]
	timestamp := param.Value["timestamp"][0]
	log.Info().Msgf("deviceId is: %v", deviceId)

	fileList := param.File["files"]
	for _, file := range fileList {
		log.Info().Msgf("file name is: %v", file.Filename)
		path := fmt.Sprintf("./static/upload/screenshots/%s", file.Filename)
		err = ctx.SaveUploadedFile(file, path)
		if err != nil {
			log.Info().Msgf("fail to save file, error is: %v", err)
		}

		var info model.ScreenReportInfo
		info.DeviceId = deviceId
		info.FilePath = path
		info.Status = 1
		var ts int64
		ts, err = strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			log.Info().Msgf("fail to parse int, error is: %v", err)
			continue
		}
		info.TimeStamp = ts
		err = repository.SaveRecordReport(&info)
		if err != nil {
			log.Info().Msgf("fail to save data, error is: %v", err)
		}
	}

	ctx.JSON(200, gin.H{
		"code":    common.Fail,
		"message": "",
	})
}

func (*ApiDevice) SetReportStatus(ctx *gin.Context) {

	type paramater struct {
		DeviceId        string `json:"device_id"`
		ReportTimeStamp int64  `json:"report_time_stamp"`
	}
	var param paramater
	err := ctx.BindJSON(&param)
	if err != nil {
		log.Info().Msgf("parameter error: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}

	log.Info().Msgf("params is: %#v", param)
	client := repository.GetRedisClient()
	client.Set(common.LastUpdateTime, param.ReportTimeStamp, 0)

	err = repository.SetRecordReport(param.DeviceId, param.ReportTimeStamp)
	if err != nil {
		log.Info().Msgf("fail to update report time, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
	}
	// 标记完成
	ctx.JSON(200, gin.H{
		"code":    common.Fail,
		"message": "",
	})
}

func (*ApiDevice) GetReportList(ctx *gin.Context) {
	redisClient := repository.GetRedisClient()
	lastUpdate, err := redisClient.Get(common.LastUpdateTime).Int64()
	log.Info().Msgf("last update is: %d", lastUpdate)
	if err != nil {
		log.Info().Msgf("fail to get last update time, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
	}
	var reportList []model.ScreenReportInfo
	reportList, err = repository.GetRecordReportList(lastUpdate)
	if err != nil {
		log.Info().Msgf("fail to get last update time, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
	}
	// 标记完成
	ctx.JSON(200, gin.H{
		"code":    common.Fail,
		"message": reportList,
	})
}

// 获取当前截屏
func (*ApiDevice) GetReport(ctx *gin.Context) {
	redisClient := repository.GetRedisClient()
	lastUpdate, err := redisClient.Get(common.LastUpdateTime).Int64()
	log.Info().Msgf("last update is: %d", lastUpdate)
	if err != nil {
		log.Info().Msgf("fail to get last update time, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
	}
	var reportList []model.ScreenReportInfo
	reportList, err = repository.GetRecordReportList(lastUpdate)
	if err != nil {
		log.Info().Msgf("fail to get last update time, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
	}
	// 标记完成
	ctx.JSON(200, gin.H{
		"code":    common.Fail,
		"message": reportList,
	})
}
