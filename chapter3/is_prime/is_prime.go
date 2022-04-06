package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
	"os"
	"strconv"
)

// isPrimeWasm was compiled from testdata/is_prime.wat
//go:embed testdata/is_prime.wasm
var isPrimeWasm []byte

func main() {
	value, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("invalid args %v: %v", os.Args[1], err)
	}

	r := wazero.NewRuntime()

	module, err := r.InstantiateModuleFromCode(isPrimeWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	isPrime := module.ExportedFunction("is_prime")
	results, err := isPrime.Call(nil, value)
	if err != nil {
		log.Fatal(err)
	}

	if results[0] == 1 {
		fmt.Printf("%d is prime!\n", value)
	} else {
		fmt.Printf("%d is NOT prime\n", value)
	}
}
