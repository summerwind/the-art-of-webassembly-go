package main

import (
	"the-art-of-webassembly-go/internal"
)

// Example_main ensures the following will work:
//
//	go build table.go
//	./table
func Example_main() {

	internal.ScrubNumbers(main)

	// Output:
	// go_table_test time=NNN
	// go_import_test time=NNN
	// wasm_table_test time=NNN
	// wasm_import_test time=NNN
}
