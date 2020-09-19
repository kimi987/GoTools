package main
import "fmt"

type C struct {
	a int
}

func main() {
	// a := C{
	// 	a: 100,
	// }

	// TestData(&a)

	// fmt.Println(a.a)

	a := 0
	if a != 0 && a != 1 {
		fmt.Println("aaa")
	}
}

func TestData(d *C) {
	d.a = 10
}

