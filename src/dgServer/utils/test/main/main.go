package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()

	fmt.Println(int(t.Weekday()))
	fmt.Println(t.Month())

	t1 := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)

	d := t.Sub(t1)
	fmt.Println(d.Hours() / 24)
}
