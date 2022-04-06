package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"log"
	"os"
	"strconv"
)

// numberStringWasm was compiled from testdata/number_string.wat
//go:embed testdata/number_string.wasm
var numberStringWasm []byte

func printString(m api.Module, strPos, strLen uint32) {
	buf, ok := m.Memory().Read(strPos, strLen)
	if !ok {
		log.Fatalf("Memory.Read(%d, %d) out of range", strPos, strLen)
	}
	fmt.Printf(">%s!\n", buf)
}

func main() {
	value, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("invalid args %v: %v", os.Args[1], err)
	}

	r := wazero.NewRuntime()

	env, err := r.NewModuleBuilder("env").
		ExportFunction("print_string", printString).
		ExportMemoryWithMax("buffer", 1, 1).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	module, err := r.InstantiateModuleFromCode(numberStringWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	toString := module.ExportedFunction("to_string")
	_, err = toString.Call(nil, value)
	if err != nil {
		log.Fatal(err)
	}
}
