package main

import (
	"github.com/tealeg/xlsx"
	"fmt"
	"os"
)


const File1 = "translation.xlsx"
const FileNew = "translation_new.xlsx"
const FileOld = "translation_old.xlsx"

var originMap = make(map[int]string)

func main() {
	xlFile_old, err := xlsx.OpenFile(FileOld)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

	sheet_old := xlFile_old.Sheet["Sheet1"]

	for _, row := range sheet_old.Rows {
		if len(row.Cells) == 0 {
			continue
		}
		key, _ := row.Cells[0].Int()

		originMap[key] = row.Cells[1].String()
	}


	xlFile, err :=  xlsx.OpenFile(File1)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}
	sheet := xlFile.Sheet["Sheet1"]

	fileNew := xlsx.NewFile()
	sheetNew, err := fileNew.AddSheet("Sheet1")

	for _, row := range sheet.Rows {
		rowNew:=sheetNew.AddRow()
		rowNew.SetHeightCM(0.5)
		if len(row.Cells) < 2 {
			continue
		}
		
		key, _ := row.Cells[0].Int()
		fmt.Println(key)
		// row.Cells[0].GetStyle().Fill.BgColor = "00000000"

		if originMap[key] != row.Cells[1].String() {


			row.Cells[1].GetStyle().Font.Color = "00FF0000"
		} else {
			// fmt.Println("originMap[key]  = " , originMap[key] )
			// 	fmt.Println("row.Cells[1].String()  = " , row.Cells[1].String())
			// row.Cells[0].GetStyle().Font.Color = "000FF000"
		}
		for _, c := range  row.Cells {
			cellNew := rowNew.AddCell()
			cellNew.Value = c.Value
			if originMap[key] != row.Cells[1].String() {


				cellNew.GetStyle().Font.Color = "00FF0000"
			} else {
				cellNew.GetStyle().Font.Color = "00000000"
			}
		}
	
	}

	err = fileNew.Save(FileNew)
    if err != nil {
        fmt.Printf(err.Error())
    }
	
}