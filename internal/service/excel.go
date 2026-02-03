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

	sheet := "Visits"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

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
	}

	for i, v := range visits {
		row := i + 2

		f.SetCellValue(sheet, "A"+strconv.Itoa(row), v.CheckInAt.Format("2006-01-02"))
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), v.Staff.FullName)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), v.Staff.PhoneNumber)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), v.Group.Name)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), v.Store.Name)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), v.CheckInAt)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), v.CheckOutAt)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), v.DistanceDeviation)
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), v.ReviewStatus)
		f.SetCellValue(sheet, "J"+strconv.Itoa(row), v.ReviewComment)
		f.SetCellValue(sheet, "K"+strconv.Itoa(row), v.Reviewer.FullName)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
