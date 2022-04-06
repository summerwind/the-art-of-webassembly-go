package main

// ExampleMain ensures the following will work:
//
//	go build globals.go
//	./globals
func ExampleMain() {
	main()

	// Output:
	// i32: 500000
	// f32: 123.012344
	// f64: 123.0123456789
}
