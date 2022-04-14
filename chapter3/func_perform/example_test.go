package main

import (
	"the-art-of-webassembly-go/internal"
)

// Example_main ensures the following will work:
//
//	go build func_perform.go
//	./func_perform
func Example_main() {

	internal.ScrubNumbers(main)

	// Output:
	// wasm_call time=NNN
	// go_call time=NNN
}
