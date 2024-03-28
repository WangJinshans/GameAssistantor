package main

import (
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

// excelize 库
// 参考: https://xuri.me/excelize/zh-hans/style.html

func GenerateExcel(templatePath string, reportPath string) {
	sheetName := "Sheet1"
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		log.Info().Msgf("fail to open template report.xlsx, error is: %v", err)
		return
	}

	// yellowWarningStyleId, _ := f.NewStyle(`{"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}],"fill":{"type":"pattern","color":["#c4bd97"],"pattern":1},"alignment":{"horizontal":"center","ident":1,"vertical":"center","wrap_text":true},"font": {"color":"#ff0000"}}`)
	blueWarningStyleId, _ := f.NewStyle(`{"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}],"fill":{"type":"pattern","color":["#b8cce4"],"pattern":1},"alignment":{"horizontal":"center","ident":1,"vertical":"center","wrap_text":true},"font": {"color":"#ff0000"}}`)

	dataFlag := true
	if dataFlag {
		f.SetCellValue(sheetName, "E5", "√")
	} else {
		// f.SetCellStyle(sheetName, "E5", "E5", yellowWarningStyleId)
		f.SetCellStyle(sheetName, "E5", "E5", blueWarningStyleId)
	}
	err = f.SaveAs(reportPath)
	if err != nil {
		log.Info().Msgf("fail to save excel file, error is: %v", err)
	}
}

func main() {
	GenerateExcel("./debug/excelize/excel.xlsx", "./debug/excelize/excel_new.xlsx")
}
