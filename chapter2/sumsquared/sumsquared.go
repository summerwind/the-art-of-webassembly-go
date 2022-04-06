package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
	"os"
	"strconv"
)

// sumsquaredWasm was compiled from testdata/sumsquared.wat
//go:embed testdata/sumsquared.wasm
var sumsquaredWasm []byte

func main() {
	val1, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("invalid args %v: %v", os.Args[1], err)
	}

	val2, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalf("invalid args %v: %v", os.Args[2], err)
	}

	r := wazero.NewRuntime()

	env, err := r.NewModuleBuilder("env").
		ExportFunction("log", func(n, factorial uint32) {
			fmt.Printf("%d! = %d\n", n, factorial)
		}).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	module, err := r.InstantiateModuleFromCode(sumsquaredWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	sumSquared := module.ExportedFunction("SumSquared")
	results, err := sumSquared.Call(nil, val1, val2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("(%d + %d) * (%d + %d) = %d\n", val1, val2, val1, val2, results[0])
}
