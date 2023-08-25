package device

import (
	"fmt"
	"game_assistantor/common"
	"game_assistantor/service"

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

func (d *ApiDevice) SendDeviceCommand(ctx *gin.Context) {

	deviceId := ctx.Param("device_id")
	log.Info().Msgf("device id is: %s", deviceId)
	if deviceId == "" {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "device id not found",
		})
		return
	}

	command := d.generateCommand(deviceId)

	err := service.SendCommand(deviceId, command)
	if err != nil {
		log.Error().Msgf("device id: %s, fail to send command, error is: %v", deviceId, err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
	}
	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
}

func (*ApiDevice) generateCommand(deviceId string) (command []byte) {
	//
	command = []byte(fmt.Sprintf("#%s#reset_all#", deviceId))
	return
}
