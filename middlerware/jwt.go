package middlerware

import (
	"game_assistantor/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, err := auth.Verify(ctx); err == nil {
			ctx.Next()
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 4001})
			ctx.Abort()
		}
	}
}
