package group_game

import "github.com/gin-gonic/gin"

var GameGroupApi ApiGameGroup

type ApiGameGroup struct {
}

// 团列表
func (g *ApiGameGroup) GroupList(ctx *gin.Context) {

}

// 构建
func (g *ApiGameGroup) RegisterGroup(ctx *gin.Context) {

}

// 更新
func (g *ApiGameGroup) UpdateGroup(ctx *gin.Context) {

}

// 延迟 带上原因 记录准确率
func (g *ApiGameGroup) PostponeGroup(ctx *gin.Context) {

}

// 申请加入 变更团的支付信息, 加入当前用户
func (g *ApiGameGroup) JoinGroup(ctx *gin.Context) {

}

// 剔掉 变更团的支付信息, 删除当前用户
func (g *ApiGameGroup) KickOutGroup(ctx *gin.Context) {

}

// 主动退出 变更团的支付信息, 删除当前用户
func (g *ApiGameGroup) LeaveGroup(ctx *gin.Context) {

}

// 解散 清除团的支付信息
func (g *ApiGameGroup) DissoluteGroup(ctx *gin.Context) {

}

// 支付证明上传
func (g *ApiGameGroup) ProveGroupFee(ctx *gin.Context) {

}

// 结果, 成功开团, 失败协商
func (g *ApiGameGroup) Result(ctx *gin.Context) {

}

// 人工审核
func (g *ApiGameGroup) ManualCheckOut(ctx *gin.Context) {

}
