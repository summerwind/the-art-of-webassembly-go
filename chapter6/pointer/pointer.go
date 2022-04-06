package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
)

// pointerWasm was compiled from testdata/pointer.wat
//go:embed testdata/pointer.wasm
var pointerWasm []byte

func main() {
	r := wazero.NewRuntime()

	env, err := r.NewModuleBuilder("env").
		ExportMemoryWithMax("buffer", 1, 1).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	module, err := r.InstantiateModuleFromCode(pointerWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	getPtr := module.ExportedFunction("get_ptr")
	results, err := getPtr.Call(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("pointer_value=%d\n", results[0])
}
