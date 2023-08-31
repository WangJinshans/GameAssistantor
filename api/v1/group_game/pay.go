package group_game

import "github.com/gin-gonic/gin"

// 支付Api
var GamePayApi ApiPayGameGroup

type ApiPayGameGroup struct {
}

// 初始化
func (g *ApiPayGameGroup) InitPayRecord(ctx *gin.Context) {
	// gorup id
}

// 用户支付, 上传orderid以及支付截图
func (g *ApiPayGameGroup) UserPay(ctx *gin.Context) {
	// group id, user id
	// 根据orderid查询支付结果
}

// 获取用户支付列表
func (g *ApiPayGameGroup) UserPayInfo(ctx *gin.Context) {
	// group id

}
