package main

import (
	"os"
)

// Example_main ensures the following will work:
//
//	go build loop.go
//	./loop 10
func Example_main() {

	// Save the old os.Args and replace with our example input.
	oldArgs := os.Args
	os.Args = []string{"loop", "10"}
	defer func() { os.Args = oldArgs }()

	main()

	// Output:
	// 1! = 1
	// 2! = 2
	// 3! = 6
	// 4! = 24
	// 5! = 120
	// 6! = 720
	// 7! = 5040
	// 8! = 40320
	// 9! = 362880
	// 10! = 3628800
	// result 10! = 3628800
}
