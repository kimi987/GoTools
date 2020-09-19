package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"resCollect/models"
)

//CreateCSVFile 创建csv文件
func CreateCSVFile(fileName string, data []*models.ResData) string {

	filePath := fmt.Sprintf("download/%s", fileName)
	_, err := os.Stat(fileName)

	if err == nil {
		err = os.Remove(fileName)

		if err != nil {
			log.Println("Remove File = ", err)
			return ""
		}
	}

	csvFile, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	// 为这个文件创建buffered writer
	bufferedWriter := bufio.NewWriterSize(
		csvFile,
		8000,
	)
	for _, v := range data {
		_, err = bufferedWriter.WriteString(
			fmt.Sprintf("%s.*,%s\n", v.ResName, v.ResPath),
		)
		if err != nil {
			log.Println(err)
		}
	}

	// 写内存buffer到硬盘
	bufferedWriter.Flush()

	return filePath
}
