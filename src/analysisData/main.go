package main

import(
	"os"
	"fmt"
	"github.com/tealeg/xlsx"
	// "strings"
	"github.com/bitly/go-simplejson"
	"time"
)

var KeyMap = make(map[string]map[int]string)

var XlsxFimeMap = make(map[string]*xlsx.File)
// func main() {
// 	data := "{\"a\":123}"
// 	json, _ := simplejson.NewJson([]byte(data))

// 	fmt.Println("json = ", json.Interface().(map[string]interface{})["a"])
// }

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("[Error] 需要参数[2], 当前[1]")
		return
	}

	fileName := os.Args[1]

	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}
	var sheet *xlsx.Sheet
	for _,v := range xlFile.Sheet {
		if v != nil {
			sheet = v;
			break
		}
	}
	// sheet := xlFile.Sheet["Sheet1"]

	if len(sheet.Rows) == 0 {
		fmt.Println("当前表为空")
        os.Exit(1)
	} 
	nameIndex := -1
	index := -1
	installTimeIndex := -1
	eventTimeIndex := -1
	for _, row := range sheet.Rows {
		if len(row.Cells) == 0 {
			continue
		}
		if index == - 1 {
			for k,cell := range row.Cells {
				if cell.Value == "Event Name" {
					nameIndex = k
				} else if cell.Value == "Event Value" {
					index = k 
				} else if cell.Value == "Install Time" {
					installTimeIndex = k
				} else if cell.Value == "Event Time" {
					eventTimeIndex = k
				}
				
			}
		} else {
			nameCell := row.Cells[nameIndex]
			cell := row.Cells[index]
			installTimeCell := row.Cells[installTimeIndex]
			eventTimeCell := row.Cells[eventTimeIndex]
			installTime, _ := installTimeCell.GetTime(false)
			eventTime ,_ := eventTimeCell.GetTime(false)
			if cell != nil {
				json, _ := simplejson.NewJson([]byte(cell.Value))
				kv := json.Interface().(map[string]interface{})
				AddRowData(nameCell.Value, installTime, eventTime, kv)
			}
		}

	}

	for k, v:= range XlsxFimeMap {
		err = v.Save(k + ".xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}


func AddRowData(eventName string, installTime, eventTime time.Time, kv map[string]interface{}) {

	// fmt.Println("installTime = ", installTime)
	if KeyMap[eventName] == nil {
		KeyMap[eventName] = make(map[int]string)

		for k,_ := range kv {
			KeyMap[eventName][len(KeyMap[eventName])] = k
		}
	}
	if XlsxFimeMap[eventName] == nil {
		XlsxFimeMap[eventName] = xlsx.NewFile()
		sheet1, _ := XlsxFimeMap[eventName].AddSheet("Sheet1")
		row := sheet1.AddRow()
		row.SetHeightCM(0.5)

		for i := 0; i < len(KeyMap[eventName]); i++ {
			cell := row.AddCell()
			cell.Value = KeyMap[eventName][i]
		}
		cell := row.AddCell()
		cell.Value = "Install Time"
		cell = row.AddCell()
		cell.Value = "Event Time"
	}

	sheet, _ := XlsxFimeMap[eventName].Sheet["Sheet1"]
	row1 := sheet.AddRow()
	row1.SetHeightCM(0.5)
	for i := 0; i < len(KeyMap[eventName]); i++ {
		key := KeyMap[eventName][i]
		val := kv[key]
		cell := row1.AddCell()
		if val == nil {
			cell.Value = ""
		} else {
			cell.Value = val.(string)
		}
		
	}
	cell := row1.AddCell()
	cell.SetDateTime(installTime)
	cell = row1.AddCell()
	cell.SetDateTime(eventTime)
}