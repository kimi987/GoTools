package main
import "io/ioutil"

func main() {
	data1, err := ioutil.ReadFile("namePart1.txt")
	data2, err := ioutil.ReadFile("namePart2.txt")
	data3, err := ioutil.ReadFile("namePart3.txt")
}

