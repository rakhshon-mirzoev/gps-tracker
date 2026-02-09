package service

import (
	"bytes"
	"strconv"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/xuri/excelize/v2"
)

type Excel interface {
	ExportReports(visits []model.VisitExcel) ([]byte, error)
}

type excelService struct{}

func NewExcelService() Excel {
	return &excelService{}
}

func (s *excelService) ExportReports(visits []model.VisitExcel) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	const sheet = "Visits"

	f.SetSheetName("Sheet1", sheet)
	f.SetActiveSheet(0)

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E7E6E6"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	dateStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 14,
	})
	if err != nil {
		return nil, err
	}

	dateTimeStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 22,
	})
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Дата визита",
		"Сотрудник",
		"Телефон",
		"Группа",
		"Магазин",
		"Check-in",
		"Check-out",
		"Отклонение (м)",
		"Статус проверки",
		"Комментарий",
		"Проверил",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, v := range visits {
		row := i + 2

		f.SetCellValue(sheet, "A"+strconv.Itoa(row), v.CheckInAt)
		f.SetCellStyle(sheet, "A"+strconv.Itoa(row), "A"+strconv.Itoa(row), dateStyle)

		f.SetCellValue(sheet, "B"+strconv.Itoa(row), v.Staff.FullName)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), v.Staff.PhoneNumber)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), v.Group.Name)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), v.Store.Name)

		f.SetCellValue(sheet, "F"+strconv.Itoa(row), v.CheckInAt)
		f.SetCellStyle(sheet, "F"+strconv.Itoa(row), "F"+strconv.Itoa(row), dateTimeStyle)

		f.SetCellValue(sheet, "G"+strconv.Itoa(row), v.CheckOutAt)
		f.SetCellStyle(sheet, "G"+strconv.Itoa(row), "G"+strconv.Itoa(row), dateTimeStyle)

		f.SetCellValue(sheet, "H"+strconv.Itoa(row), v.DistanceDeviation)
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), v.ReviewStatus)
		f.SetCellValue(sheet, "J"+strconv.Itoa(row), v.ReviewComment)
		f.SetCellValue(sheet, "K"+strconv.Itoa(row), v.Reviewer.FullName)
	}

	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheet, colName, colName, 20)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
