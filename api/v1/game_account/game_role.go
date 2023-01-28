package game_account

import (
	"game_assistantor/common"
	"game_assistantor/repository"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var GameRoleApi ApiGameRole

type ApiGameRole struct {
}

func (*ApiGameRole) GetAccountInfo(ctx *gin.Context) {
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

func (*ApiGameRole) GetAccountRoleList(ctx *gin.Context) {
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

func (*ApiGameRole) UpdateRoleInfo(ctx *gin.Context) {
	type req struct {
		UserId  string `json:"user_id"`
		OrderId string `json:"order_id"`
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

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
	return
}

func (*ApiGameRole) UpdateAccountInfo(ctx *gin.Context) {
	type req struct {
		UserId  string `json:"user_id"`
		OrderId string `json:"order_id"`
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

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
	return
}

// 添加账号
func (*ApiGameRole) AddAccount(ctx *gin.Context) {
	type req struct {
		UserId  string `json:"user_id"`
		OrderId string `json:"order_id"`
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

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
	return
}

// 添加角色
func (*ApiGameRole) AddAccountRole(ctx *gin.Context) {
	type req struct {
		UserId  string `json:"user_id"`
		OrderId string `json:"order_id"`
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

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
	return
}
