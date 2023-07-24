package group_game

import "github.com/gin-gonic/gin"

// 信誉

var GameReputationApi ApiReputation

type ApiReputation struct {
}

// 完成单个团, 信誉增加(开团次数, 成功次数, 失败次数)
func (g *ApiReputation) FinishSingleGroupGame(ctx *gin.Context) {

}

// 开始
func (g *ApiReputation) StartSingleGroupGame(ctx *gin.Context) {

}
