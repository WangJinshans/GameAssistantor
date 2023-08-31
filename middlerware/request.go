package middlerware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func RequestInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		c.Next()
		end := time.Since(start)
		log.Info().Msgf("request url is: %s, method is: %s, time cost is: %v", c.Request.RequestURI, method, end)
	}
}
