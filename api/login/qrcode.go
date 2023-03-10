package login

import (
	"encoding/base64"
	"fmt"
	"game_assistantor/auth"
	"game_assistantor/common"
	"game_assistantor/global"
	"game_assistantor/model"
	"game_assistantor/repository"
	"game_assistantor/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	qrcode "github.com/skip2/go-qrcode"
	"path/filepath"
	"time"
)

var statusList = []int{
	common.InitStatus,
	common.ScanStatus,
	common.UnConfirmStatus,
	common.CancelStatus,
	common.SuccessStatus,
	common.FailStatus,
	common.InValidStatus,
}

// 状态校验
func IsValidStatus(status int) bool {
	for _, item := range statusList {
		if status == item {
			return true
		}
	}
	return false
}

// 初始化二维码
func InitQrCode(ctx *gin.Context) () {
	qrcodeId := fmt.Sprintf("%s_%d", utils.GetUuidString(), time.Now().Unix())
	qrcodeName := fmt.Sprintf("%s.png", qrcodeId)
	content, err := repository.GenerateQrcodeContent(qrcodeId, qrcodeName)
	if err != nil {
		log.Info().Msgf("fail to init qr code, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}

	var baseDir string
	baseDir, err = global.GetExecutablePath()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
	}
	qrcodeName = filepath.Join(baseDir, "static", "upload", "qrcode", qrcodeName) // file path
	err = qrcode.WriteFile(string(content), qrcode.Medium, 256, qrcodeName)
	if err != nil {
		log.Info().Msgf("fail to init qr code, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}

	qrCode := model.QrCodeInfo{
		QrCodeId:  qrcodeId,
		IsValid:   "true",
		ImagePath: "",
		TimeStamp: 0,
		Status:    common.InitStatus,
		Content:   "",
	}
	err = repository.SaveQrCode(qrCode) // 创建
	if err != nil {
		log.Info().Msgf("fail to init qr code, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}
	message := base64.StdEncoding.EncodeToString(content) // base64
	ctx.JSON(200, gin.H{
		"code":      common.Fail,
		"message":   message,
		"qrcode_id": qrcodeId,
	})
	return
}

// 查询二维码状态
func QueryQrCode(ctx *gin.Context) () {
	type req struct {
		QrcodeId string `json:"qrcode_id"`
	}
	var parameter req
	err := ctx.BindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "no qrcode id found",
		})
	}

	info, err := repository.GetQrCode(parameter.QrcodeId)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
		"status":  info.Status,
	})
	return
}

// 扫描二维码状态
func ScanQrCode(ctx *gin.Context) () {
	type req struct {
		Token    string `json:"token"`
		QrCodeId string `json:"qrcode_id"`
	}
	var parameter req
	err := ctx.BindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "no qrcode id found",
		})
	}

	claim, err := auth.VerifyToken(parameter.Token)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "token error",
		})
		return
	}

	userId := claim.UserId
	qrCode := model.QrCodeInfo{
		QrCodeId: parameter.QrCodeId,
		IsValid:  "true",
		Status:   common.ScanStatus,
		UserId:   userId,
	}

	err = repository.SaveQrCode(qrCode)
	if err != nil {
		log.Info().Msgf("fail to init qr code, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}
	return
}

// 更新二维码状态
func SetQrCodeStatus(ctx *gin.Context) () {
	type req struct {
		Token    string `json:"token"`
		QrCodeId string `json:"qrcode_id"`
		Status   int    `json:"status"`
	}
	var parameter req
	err := ctx.BindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "no qrcode id found",
		})
	}

	claim, err := auth.VerifyToken(parameter.Token)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "token error",
		})
		return
	}

	if !IsValidStatus(parameter.Status) {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "invalid status",
		})
		return
	}

	userId := claim.UserId
	qrCode := model.QrCodeInfo{
		QrCodeId: parameter.QrCodeId,
		IsValid:  "true",
		Status:   parameter.Status,
		UserId:   userId,
	}

	err = repository.SaveQrCode(qrCode)
	if err != nil {
		log.Info().Msgf("fail to init qr code, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}
	return
}
