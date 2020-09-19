package main

import (
	// "time"
	"fmt"
	// "math/rand"
)

func main() {
	// nowTime := time.Now()
	// fmt.Println(nowTime.UnixNano())
	// var sum = 0
	// for i := 0; i < 500000; i++ {
	// 	if i > 0 && i < 10000 || i > 10000 && i < 20000 {
	// 		sum += i
	// 	}
	// }
	// fmt.Println(time.Now().Sub(nowTime))
	// fmt.Println(math.Pow(2, 2))
	// a := strings.Split("","")
	// fmt.Println(a)
	// fmt.Printf("%b", GetByte(65536) )
	
	// s := []string      {0: "no error", 1: "Eio", 2: "invalid argument"}
	// fmt.Println(s)
	// maxInt :=  int(^uint(0) >> 1)
	// fmt.Printf("%b", 5 ^ maxInt)

	list := make([]string, 0 ,32)

	fmt.Println(len(list))
	list = append(list,"112")
	fmt.Println(len(list))
}


func GetByte(id int)[4]byte {
	return [4]byte{byte(id >>24), byte(id >> 16&255),byte(id>>8&255), byte(id&255)}
}