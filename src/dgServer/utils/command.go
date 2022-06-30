package utils

import (
	"fmt"
	"os/exec"
	"time"
)

func ExecFile(filename string) bool {

	fmt.Println("filename = ", filename)

	// c := exec.Command("start", filename)
	c := exec.Command(filename)

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)

		return false
	}

	return true

}

//GetNowTimeStringHour 获取当前的小时数据
func GetNowTimeStringHour() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
