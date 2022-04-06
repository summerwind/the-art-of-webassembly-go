package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
)

// globalsWasm was compiled from testdata/globals.wat
//go:embed testdata/globals.wasm
var globalsWasm []byte

func log_i32(value int32) {
	fmt.Printf("i32: %v\n", value)
}
func log_f32(value float32) {
	fmt.Printf("f32: %v\n", value)
}
func log_f64(value float64) {
	fmt.Printf("f64: %v\n", value)
}

func main() {
	r := wazero.NewRuntime()

	env, err := r.NewModuleBuilder("env").
		ExportGlobalI32("import_i32", 500_000). // 5_000_000 is out of range of i32!
		ExportGlobalF32("import_f32", 123.0123456789).
		ExportGlobalF64("import_f64", 123.0123456789).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	goM, err := r.NewModuleBuilder("go").
		ExportFunction("log_i32", log_i32).
		ExportFunction("log_f32", log_f32).
		ExportFunction("log_f64", log_f64).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer goM.Close()

	module, err := r.InstantiateModuleFromCode(globalsWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	test := module.ExportedFunction("globaltest")
	_, err = test.Call(nil)
	if err != nil {
		log.Fatal(err)
	}
}
