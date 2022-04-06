package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"log"
)

// helloworldWasm was compiled from testdata/helloworld.wat
//go:embed testdata/helloworld.wasm
var helloworldWasm []byte

var startStringOffset = int32(100)

func printString(m api.Module, strLen uint32) {
	buf, ok := m.Memory().Read(uint32(startStringOffset), strLen)
	if !ok {
		panic("out of range reading string")
	}
	fmt.Printf("%s\n", buf)
}

func main() {
	r := wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfigInterpreter())

	envImports, err := r.NewModuleBuilder("env").
		ExportFunction("print_string", printString).
		ExportMemory("buffer", 1).
		ExportGlobalI32("start_string", startStringOffset).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer envImports.Close()

	module, err := r.InstantiateModuleFromCode(helloworldWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	helloworld := module.ExportedFunction("helloworld")
	_, err = helloworld.Call(nil)
	if err != nil {
		log.Fatal(err)
	}
}
