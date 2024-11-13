package excelGenerator

import (
	"bytes"

	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/xuri/excelize/v2"
)

func GenerateExcelReport(reportData types.ReportData, buf *bytes.Buffer) error {

	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", reportData.Title)
	rows_written := 1
	if len(reportData.CommonDataHeads) > 0 {
		rows_written++
		style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})

		for index1, value1 := range reportData.ConvertTo2D(reportData.CommonDataHeads, 3) {

			for index2, value2 := range value1 {

				cellName1, _ := excelize.CoordinatesToCellName((index2*2)+1, 1+rows_written)
				f.SetCellValue("Sheet1", cellName1, value2)
				f.SetCellStyle("Sheet1", cellName1, cellName1, style)
				cellName2, _ := excelize.CoordinatesToCellName((index2*2)+2, 1+rows_written)
				f.SetCellValue("Sheet1", cellName2, reportData.CommonDataDetails[(index1*3)+index2])

			}
			rows_written++
		}

	}

	rows_written++

	colHeadingRow := 1 + rows_written
	//set column headings
	for i := 0; i < len(reportData.ColumnHeadings); i++ {

		cellName, _ := excelize.CoordinatesToCellName(i+1, 1+rows_written)
		f.SetCellValue("Sheet1", cellName, reportData.ColumnHeadings[i])

	}

	rows_written++
	//set detail rows

	for i := 0; i < len(reportData.DataRows); i++ {

		for j := 0; j < len(reportData.DataRows[i]); j++ {
			cellName, _ := excelize.CoordinatesToCellName(j+1, 1+rows_written)

			f.SetCellValue("Sheet1", cellName, reportData.DataRows[i][j])

		}
		rows_written++
	}

	//styling cells

	style, err := f.NewStyle(&excelize.Style{Fill: excelize.Fill{Type: "pattern", Color: []string{"#DFEBF6"}, Pattern: 1}, Font: &excelize.Font{Bold: true, Underline: "double"}})
	if err != nil {
		return err
	}
	f.SetCellStyle("Sheet1", "A1", "E1", style)

	style, err = f.NewStyle(&excelize.Style{Fill: excelize.Fill{Type: "pattern", Color: []string{"#FDFE7E"}, Pattern: 1}, Font: &excelize.Font{Bold: true}, Border: []excelize.Border{{Type: "top", Style: 2, Color: "000000"}, {Type: "left", Style: 2, Color: "000000"}, {Type: "right", Style: 2, Color: "000000"}, {Type: "bottom", Style: 2, Color: "000000"}}})

	if err != nil {
		return err
	}

	startcellName, _ := excelize.CoordinatesToCellName(1, colHeadingRow)
	endcellName, _ := excelize.CoordinatesToCellName(len(reportData.ColumnHeadings), colHeadingRow)
	f.SetCellStyle("Sheet1", startcellName, endcellName, style)

	//write footer if required

	if len(reportData.FooterDataHeads) > 0 {

		rows_written++

		cellName1, _ := excelize.CoordinatesToCellName(1, 1+rows_written)

		style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Underline: "single"}})

		f.SetCellValue("Sheet1", cellName1, "Summary")
		f.SetCellStyle("Sheet1", cellName1, cellName1, style)

		style, _ = f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})

		rows_written++
		for index1, value1 := range reportData.ConvertTo2D(reportData.FooterDataHeads, 3) {

			for index2, value2 := range value1 {

				cellName1, _ := excelize.CoordinatesToCellName((index2*2)+1, 1+rows_written)
				f.SetCellValue("Sheet1", cellName1, value2)
				f.SetCellStyle("Sheet1", cellName1, cellName1, style)
				cellName2, _ := excelize.CoordinatesToCellName((index2*2)+2, 1+rows_written)
				f.SetCellValue("Sheet1", cellName2, reportData.FooterDataDetails[(index1*3)+index2])

			}
			rows_written++
		}

	}

	f.Write(buf)

	return nil

}
