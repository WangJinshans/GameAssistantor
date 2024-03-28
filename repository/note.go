package repository

import "game_assistantor/model"

// 保存笔记
func SaveNote(info model.StockNote) (err error) {
	err = GetEngine().Save(&info).Error
	return
}

// 获取列表
func GetNoteList(userId string) (err error) {
	var noteList []model.StockNote
	err = GetEngine().Model(&model.StockNote{}).Where("user_id = ?", userId).Find(&noteList).Error
	return
}

// 更新笔记
func UpdateNoteInfo(info model.StockNote) (err error) {
	err = GetEngine().Save(&info).Error
	return
}
