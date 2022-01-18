package main

import (
	"fmt"
	"math"
	"os"
	"reflect"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

func log_i32(ctx *wasm.HostFunctionCallContext, value int32) {
	fmt.Printf("i32: %v\n", value)
}
func log_f32(ctx *wasm.HostFunctionCallContext, value float32) {
	fmt.Printf("f32: %v\n", value)
}
func log_f64(ctx *wasm.HostFunctionCallContext, value float64) {
	fmt.Printf("f64: %v\n", value)
}

func main() {
	buf, err := os.ReadFile("globals.wasm")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read wasm file: %v", err)
		os.Exit(1)
	}

	mod, err := wasm.DecodeModule(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode module: %v", err)
		os.Exit(1)
	}

	store := wasm.NewStore(wazeroir.NewEngine())

	err = store.AddHostFunction("go", "log_i32", reflect.ValueOf(log_i32))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}
	err = store.AddHostFunction("go", "log_f32", reflect.ValueOf(log_f32))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}
	err = store.AddHostFunction("go", "log_f64", reflect.ValueOf(log_f64))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	store.ModuleInstances["env"] = &wasm.ModuleInstance{Exports: map[string]*wasm.ExportInstance{}}

	err = store.AddGlobal("env", "import_i32", uint64(5000000000), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}
	err = store.AddGlobal("env", "import_f32", uint64(math.Float32bits(123.0123456789)), wasm.ValueTypeF32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}
	err = store.AddGlobal("env", "import_f64", uint64(math.Float64bits(123.0123456789)), wasm.ValueTypeF64, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	_, _, err = store.CallFunction("", "globaltest")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
}
