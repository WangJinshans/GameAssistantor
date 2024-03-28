package auth

import (
	"game_assistantor/global"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// casbin 鉴权
func VerifyAccess(ctx *gin.Context) {
	var result bool

	type reqinfo struct {
		UserId string `json:"user_id"`
		Path   string `json:"path"`
		Method string `json:"method"`
	}

	var info reqinfo
	err := ctx.BindJSON(&info)
	if err != nil {
		log.Info().Msg("headers invalid")
		ctx.JSON(200, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}

	enforcer := global.GetEnforcer()
	result, err = enforcer.Enforce(info.UserId, info.Path, info.Method)
	// Enforce 会验证角色的相关的权限
	// HasPermissionForUser 只验证用户是否有权限
	// res = Enforcer.HasPermissionForUser(userName,p,m)
	if err != nil {
		log.Info().Msgf("user id: %s, has no permission for path: %s, error is: %v", info.UserId, info.Path, err)
		ctx.JSON(200, gin.H{
			"message": "Unauthorized",
			"data":    "",
		})
		ctx.Abort()
		return
	}
	if !result {
		log.Info().Msgf("user id: %s, path: %s, permission check failed", info.UserId, info.Path)
		ctx.JSON(200, gin.H{
			"message": "Unauthorized",
			"data":    "",
		})
		ctx.Abort()
		return
	}

	Verify(ctx)
}
