package game_account

import (
	"game_assistantor/common"
	"game_assistantor/model"
	"game_assistantor/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var GameRoleApi ApiGameRole

type ApiGameRole struct {
}

func (*ApiGameRole) GetAccountInfo(ctx *gin.Context) {
	type req struct {
		AccountId string `json:"account_id"`
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

	user, err := repository.GameRoleRepos.GetAccountInfo(parameter.AccountId)
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
		AccountId string `json:"account_id"`
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
	var roleList []model.GameRole
	roleList, err = repository.GameRoleRepos.GetAccountRoleList(parameter.AccountId)
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
		"message": roleList,
	})
	return
}

func (*ApiGameRole) UpdateAccountInfo(ctx *gin.Context) {
	type req struct {
		AccountId string `json:"account_id"`
	}
	var parameter req
	err := ctx.BindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "account not found",
		})
		return
	}
	var account model.GameAccount
	account.AccountId = parameter.AccountId
	err = repository.GameRoleRepos.SaveAccount(account)
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
		"message": "",
	})
	return
}

// 添加账号
func (*ApiGameRole) AddAccount(ctx *gin.Context) {

	validate := validator.New()

	type req struct {
		AccountId   string `json:"account_id" validate:"required"`
		AccountPwd  string `json:"account_pwd" validate:"required"`
		AccountType string `json:"account_type" validate:"required"`
	}
	var parameter req
	err := ctx.ShouldBindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "user id not found",
		})
		return
	}

	err = validate.Struct(&parameter)
	if err != nil {

	}
	var account model.GameAccount
	account.AccountId = parameter.AccountId
	account.AccountType = parameter.AccountType
	account.AccountPwd = parameter.AccountPwd
	err = repository.GameRoleRepos.SaveAccount(account)
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
		"message": "",
	})
	return
}

// 添加角色
func (*ApiGameRole) AddAccountRole(ctx *gin.Context) {
	type req struct {
		AccountId    string `json:"account_id"`
		RoleId       string `json:"role_id"`
		RoleName     string `json:"role_name"`
		SummonerType string `json:"summoner_type"` // 角色类型
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

func (*ApiGameRole) UpdateRoleInfo(ctx *gin.Context) {
	type req struct {
		RoleId string `json:"role_id"`
	}
	var parameter req
	err := ctx.BindJSON(&parameter)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": "role id not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
	return
}
