package middlerware

import (
	"game_assistantor/global"
	"game_assistantor/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var result bool
		userId, err := c.Cookie("user_id")
		if userId == "" {
			log.Info().Msg("headers invalid")
			c.JSON(200, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		p := c.Request.URL.Path
		m := c.Request.Method

		enforcer := global.GetEnforcer()
		result, err = enforcer.Enforce(userId, p, m)
		// Enforce 会验证角色的相关的权限
		// HasPermissionForUser 只验证用户是否有权限
		//res = Enforcer.HasPermissionForUser(userName,p,m)
		if err != nil {
			log.Info().Msgf("user id: %s, has no permission for path: %s, error is: %v", userId, p, err)
			c.JSON(200, gin.H{
				"message": "Unauthorized",
				"data":    "",
			})
			c.Abort()
			return
		}
		if !result {
			log.Info().Msgf("user id: %s, path: %s, permission check failed", userId, p)
			c.JSON(200, gin.H{
				"message": "Unauthorized",
				"data":    "",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RomoteAuthMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		dataMap := make(map[string]interface{})
		// 地址固定, 方法固定, 参数变动
		dataMap["token"] = ""
		dataMap["user_id"] = "" // user id for casbin
		dataMap["path"] = ""    // 路径 for casbin
		dataMap["method"] = ""  // 方法 for casbin
		utils.DoRequest("", "POST", dataMap)
	}
}
