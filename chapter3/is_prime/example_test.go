package main

import (
	"os"
)

// Example_main ensures the following will work:
//
//	go build is_prime.go
//	./is_prime 7
func Example_main() {

	// Save the old os.Args and replace with our example input.
	oldArgs := os.Args
	os.Args = []string{"is_prime", "7"}
	defer func() { os.Args = oldArgs }()

	main()

	// Output:
	// 7 is prime!
}
