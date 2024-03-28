package note

import (
	"game_assistantor/common"
	"game_assistantor/model"
	"game_assistantor/repository"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var NoteApi ApiNote

type ApiNote struct {
}

// 单个笔记
func (*ApiNote) GetNote(ctx *gin.Context) {
	// var noteId string
	userId := ctx.Param("user_id")
	noteId := ctx.Param("note_id")

	log.Info().Msgf("user id is: %s, note id is: %s", userId, noteId)
}

// 笔记列表
func (*ApiNote) GetNoteList(ctx *gin.Context) {
	// var userId string
	userId := ctx.Param("user_id")
	log.Info().Msgf("user id is: %s", userId)

}

// 创建笔记
func (*ApiNote) CreateNote(ctx *gin.Context) {
	var info model.StockNote
	err := ctx.ShouldBind(&info)
	if err != nil {
		log.Info().Msgf("fail to create note, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}

	err = repository.SaveNote(info)
	if err != nil {
		log.Info().Msgf("fail to create note, error is: %v", err)
		ctx.JSON(200, gin.H{
			"code":    common.Fail,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    common.Success,
		"message": "",
	})
}

// 删除笔记
func (*ApiNote) DeleteNote(ctx *gin.Context) {
	// var noteId string

}

// 更新笔记
func (*ApiNote) UpdateNote(ctx *gin.Context) {
	// var info model.StockNote
}
