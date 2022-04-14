package main

import (
	"os"
)

// Example_main ensures the following will work:
//
//	go build number_string.go
//	./number_string 4103
func Example_main() {

	// Save the old os.Args and replace with our example input.
	oldArgs := os.Args
	os.Args = []string{"number_string", "4103"}
	defer func() { os.Args = oldArgs }()

	main()

	// Output:
	// >            4103!
	// >          0x1007!
	// > 0000 0000 0000 0000 0001 0000 0000 0111!
}
