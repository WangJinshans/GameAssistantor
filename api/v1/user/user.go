package user

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
	"time"
)

var UserApi ApiUser

type ApiUser struct {
}

func (*ApiUser) GetUsersInfo(ctx *gin.Context) {
	type req struct {
		UserId string `json:"user_id"`
	}
	var parameter req
	err := ctx.BindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "user id not found",
		})
		return
	}

	user, err := repository.UserRepos.GetUserList()
	if err != nil {
		log.Info().Msgf("fail to get user, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "user not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": user,
	})
	return
}

func (*ApiUser) GetUserInfo(ctx *gin.Context) {
	type req struct {
		UserId string `json:"user_id"`
	}
	var parameter req
	err := ctx.BindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "user id not found",
		})
		return
	}

	user, err := repository.UserRepos.GetUser(parameter.UserId)
	if err != nil {
		log.Info().Msgf("fail to get user, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "user not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": user,
	})
	return
}

func (*ApiUser) AddUserInfo(context *gin.Context) {

	var user model.User
	err := context.BindJSON(&user)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "password error",
		})
		return
	}

	contextId := context.GetHeader("ctx_id")
	log.Info().Msgf("encrypt password is %s", user.PassWord)
	data, err := base64.StdEncoding.DecodeString(user.PassWord)
	if err != nil {
		log.Error().Msgf("base64 decode error: %v", err.Error())
		return
	}

	privateKey, err := global.GetPrivateKey(contextId)
	if err != nil {
		log.Error().Msgf("Decrypt error: %v", err)
		return
	}

	passWord, err := utils.RsaDecrypt(data, privateKey)
	if err != nil {
		log.Error().Msgf("rsa decrypt error: %v", err.Error())
		return
	}
	global.DeleteKey(contextId)
	userId := fmt.Sprintf("user_%d", time.Now().Unix())
	user.UserId = userId
	user.PassWord = string(passWord)
	repository.UserRepos.SaveUserPassword(&user)

	var token string
	token, err = auth.GenerateToken(userId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "failed",
			"token":   "",
			"user_id": "",
		})
		return
	}
	err = repository.SaveToken(token, userId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "failed",
			"token":   "",
			"user_id": "",
		})
		return
	}
	context.JSON(200, gin.H{
		"message": userId,
		"token":   token,
		"user_id": userId,
	})
}

func (*ApiUser) UpdateUserInfo(context *gin.Context) {

	var user model.User
	err := context.BindJSON(&user)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "password error",
		})
		return
	}

	repository.UserRepos.UpdateUser(&user)
	context.JSON(200, gin.H{
		"message": "",
	})
}

func (*ApiUser) UpdateUserPassword(context *gin.Context) {

	var user model.User
	err := context.BindJSON(&user)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "parameter error",
		})
		return
	}

	err = repository.UserRepos.SaveUserPassword(&user)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "password save error",
		})
		return
	}
	context.JSON(200, gin.H{
		"message": "",
	})
}
