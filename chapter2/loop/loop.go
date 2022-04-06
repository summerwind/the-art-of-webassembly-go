package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
	"os"
	"strconv"
)

// loopWasm was compiled from testdata/loop.wat
//go:embed testdata/loop.wasm
var loopWasm []byte

func main() {
	n, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("invalid args %v: %v", os.Args[1], err)
	}
	if n > 12 {
		log.Fatal("factorials greater than 12 are too large for a 32-bit integer")
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

	module, err := r.InstantiateModuleFromCode(loopWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	loop := module.ExportedFunction("loop_test")
	results, err := loop.Call(nil, n)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("result %d! = %d\n", n, results[0])
}
