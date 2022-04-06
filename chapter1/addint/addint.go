package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
	"os"
	"strconv"
)

// addintWasm was compiled from testdata/addint.wat
//go:embed testdata/addint.wasm
var addintWasm []byte

func main() {
	r := wazero.NewRuntime()

	val1, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("invalid args %v: %v", os.Args[1], err)
	}

	val2, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalf("invalid args %v: %v", os.Args[2], err)
	}

	module, err := r.InstantiateModuleFromCode(addintWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	addInt := module.ExportedFunction("AddInt")
	results, err := addInt.Call(nil, val1, val2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d + %d = %d\n", val1, val2, results[0])
}
