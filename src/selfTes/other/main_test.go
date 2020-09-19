	package main
	import "testing"

	func TestMain(test *testing.T) {
		main()
	}


	func BenchmarkMain(bmtest *testing.B) {
		for cnt := 0; cnt < bmtest.N; cnt++ {
			main()
		} 
	}