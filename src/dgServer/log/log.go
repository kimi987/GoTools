package log

import (
	"dgServer/conf"
	"fmt"
	"log"
	"os"
	"time"
)

var logger *log.Logger

func init() {
	nowString := time.Now().Unix()
	path := fmt.Sprintf("%s/%d.%s", conf.Config.LogPath, nowString, "log")
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)
}

//Debug 输出
func Debug(format string, values ...interface{}) {
	value := fmt.Sprintf(format, values...)
	fmt.Println("value = ", value)
	logger.Println(value)
}

//Error 错误
func Error(format string, values ...interface{}) {
	value := fmt.Sprintf(format, values...)
	fmt.Println("[ERROR] ", value)
	logger.Println("[ERROR] ", value)
}
