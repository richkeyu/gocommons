package excel

import (
	"context"
	"fmt"

	"github.com/richkeyu/gocommons/plog"

	"github.com/xuri/excelize/v2"
)

// ExportSliceToExcel 导出数组到本地excel文件
func ExportSliceToExcel(ctx context.Context, sheet []string, data [][][]string, path string, extFunc func(file *excelize.File)) error {
	if len(sheet) != len(data) {
		return fmt.Errorf("ExportSliceToExcel sheet len != data len")
	}
	f := excelize.NewFile()
	for i, sd := range data {
		s := sheet[i]
		f.NewSheet(s)
		for i, line := range sd {
			for i2, value := range line {
				name, err := excelize.CoordinatesToCellName(i2+1, i+1)
				if err != nil {
					plog.Errorf(ctx, "gen excel file get name fail: %s col: %d row: %d", err, i2, i)
					continue
				}
				err = f.SetCellValue(s, name, value)
				if err != nil {
					plog.Errorf(ctx, "gen excel file set value fail: %s col: %d row: %d", err, i2, i)
					continue
				}
			}
		}
	}
	f.DeleteSheet("Sheet1")

	// 扩展处理
	if extFunc != nil {
		extFunc(f)
	}

	err := f.SaveAs(path)
	if err != nil {
		return err
	}
	return nil
}

// SliceToExcel 导出内存中
func SliceToExcel(ctx context.Context, sheet []string, data [][][]string,
	extFunc func(file *excelize.File)) (*excelize.File, error) {
	if len(sheet) != len(data) {
		return nil, fmt.Errorf("ExportSliceToExcel sheet len != data len")
	}
	f := excelize.NewFile()
	for i, sd := range data {
		s := sheet[i]
		f.NewSheet(s)
		for i, line := range sd {
			for i2, value := range line {
				name, err := excelize.CoordinatesToCellName(i2+1, i+1)
				if err != nil {
					plog.Errorf(ctx, "gen excel file get name fail: %s col: %d row: %d", err, i2, i)
					continue
				}
				err = f.SetCellValue(s, name, value)
				if err != nil {
					plog.Errorf(ctx, "gen excel file set value fail: %s col: %d row: %d", err, i2, i)
					continue
				}
			}
		}
	}
	f.DeleteSheet("Sheet1")

	// 扩展处理
	if extFunc != nil {
		extFunc(f)
	}

	return f, nil
}
