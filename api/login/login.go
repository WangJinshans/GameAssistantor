package login

import (
	"encoding/base64"
	"fmt"
	"game_assistantor/api/v1/user"
	"game_assistantor/auth"
	"game_assistantor/common"
	"game_assistantor/global"
	"game_assistantor/model"
	"game_assistantor/repository"
	"game_assistantor/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Response struct {
	Message string `json:"message,omitempty"`
	CtxId   string `json:"ctx_id"`
	UserId  string `json:"user_id"`
}

type LoginResponse struct {
	Message string `json:"message,omitempty"`
	UserId  string `json:"user_id"`
	Token   string `json:"token"`
}

// @获取rsa公钥
// @Description get rsa public key
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /get_public_key [get]
func GetPublicKey(c *gin.Context) {
	publicKeyStr, privateKeyStr, err := utils.GenRsaKey(1024)
	if err != nil {
		log.Error().Msg("密钥文件生成失败！")
		c.JSON(200, gin.H{
			"message": "server error",
		})
		return
	}
	key := fmt.Sprintf("%s%d", uuid.New().String(), time.Now().Unix())
	global.AddKeyValue(key, privateKeyStr, publicKeyStr)
	c.JSON(200, gin.H{
		"message": publicKeyStr,
		"ctx_id":  key,
	})
}

// @登录
// @Description post 登录
// @Produce json
// @Param user_name formData string true "用户名"
// @Param pass_word formData string true "密码"
// @Success 200 {object} LoginResponse
// @Router /login [post]
func Login(context *gin.Context) {
	var user model.User
	err := context.BindJSON(&user)
	if err != nil {
		log.Info().Msg(err.Error())
		context.JSON(200, gin.H{
			"message": "password error",
		})
		return
	}
	contextId := context.GetHeader("ctx_id")
	privateKey, err := global.GetPrivateKey(contextId)
	if err != nil {
		log.Error().Msgf("Decrypt error: %v", err)
		return
	}

	var data []byte
	data, err = base64.StdEncoding.DecodeString(user.PassWord)
	if err != nil {
		log.Error().Msgf("base64 error is: %v", err.Error())
		return
	}

	passWord, err := utils.RsaDecrypt(data, privateKey)
	if err != nil {
		log.Error().Msgf("rsa decrypt error: %v", err.Error())
		return
	}
	u, err := repository.UserRepos.GetUser(user.UserId)
	if err != nil {
		log.Error().Msgf("error is: %v", err.Error())
		return
	}
	log.Info().Msgf("password is: %s, data base password is: %s", passWord, u.PassWord)
	err = bcrypt.CompareHashAndPassword([]byte(u.PassWord), passWord) //验证（对比）
	if err != nil {
		log.Error().Msgf("error is: %v", err.Error())
		return
	}
	var token string
	token, err = auth.GenerateToken(u.UserId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "failed",
			"token":   "",
			"user_id": "",
		})
		return
	}
	global.DeleteKey(contextId)
	err = repository.SaveToken(token, u.UserId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "failed",
			"token":   "",
			"user_id": "",
		})
		return
	}
	context.JSON(200, gin.H{
		"message": "success",
		"token":   token,
		"user_id": u.UserId,
	})
}

func Register(context *gin.Context) {
	user.UserApi.AddUserInfo(context)
}

func RefreshToken(ctx *gin.Context) {
	token, userId, err := auth.Refresh(ctx)
	if err != nil {
		log.Info().Msgf("refresh token error: %v", err)
		ctx.Abort()
		return
	}
	if token != "" {
		ctx.JSON(200, gin.H{
			"code":    common.Success,
			"message": "",
		})
		return
	}
	err = repository.SaveToken(token, userId)
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "failed",
			"token":   "",
			"user_id": "",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "empty token",
	})
	return
}

func Logout(ctx *gin.Context) {
	var user model.User
	err := ctx.BindJSON(&user)
	if err != nil {
		log.Info().Msgf("logout error: %v", err.Error())
		ctx.JSON(200, gin.H{
			"message": "password error",
		})
		return
	}

	repository.RemoveToken(user.UserId)
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
	return
}
