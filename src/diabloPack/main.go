package main

import(
	"os"
	"os/exec"
	"fmt"
	"time"
)

var command = "./package.sh 1>logPack.txt"
var cmd = exec.Command("/bin/bash", "-c", command)
var tick = time.Tick(30 * time.Second)

func main() {
	CheckFileByTime()
}

func CheckFileByTime() {
	for {
		select {
		case <-tick:
			CheckFile()
		}
	}
}

func CheckFile() {

	_, err := os.Stat("Temp.zip")
	if err != nil {
		//文件不存在
		// fmt.Println("err = ", err)
		return
	}
	err = cmd.Run()
    if err != nil {
		// fmt.Println("Execute Command failed:" + err.Error())
		cmd = exec.Command("/bin/bash", "-c", command)
        return
	}
	fmt.Println("Execute Command !")

}
