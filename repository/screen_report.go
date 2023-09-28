package repository

import "game_assistantor/model"

func SaveRecordReport(info *model.ScreenReportInfo) (err error) {
	err = engine.Save(info).Error
	return
}

func GetRecordReportList(lastUpdate int64) (reportList []model.ScreenReportInfo, err error) {
	err = engine.Debug().Model(&model.ScreenReportInfo{}).Where("time_stamp > ? and status = 1", lastUpdate).Find(&reportList).Error
	return
}

func SetRecordReport(deviceId string, reportTimeStamp int64) (err error) {
	err = engine.Model(&model.ScreenReportInfo{}).Where("time_stamp = ? and device_id = ?", reportTimeStamp, deviceId).Update("status", "2").Error
	return
}
