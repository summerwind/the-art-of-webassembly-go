package main

import (
	"the-art-of-webassembly-go/internal"
)

// ExampleMain ensures the following will work:
//
//	go build table.go
//	./table
func ExampleMain() {

	internal.ScrubNumbers(main)

	// Output:
	// go_table_test time=NNN
	// go_import_test time=NNN
	// wasm_table_test time=NNN
	// wasm_import_test time=NNN
}
