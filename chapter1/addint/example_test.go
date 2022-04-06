package main

import (
	"os"
)

// ExampleMain ensures the following will work:
//
//	go build addint.go
//	./addint 7 9
func ExampleMain() {

	// Save the old os.Args and replace with our example input.
	oldArgs := os.Args
	os.Args = []string{"addint", "7", "9"}
	defer func() { os.Args = oldArgs }()

	main()

	// Output:
	// 7 + 9 = 16
}
