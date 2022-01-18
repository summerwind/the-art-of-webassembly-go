package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

var i int32 = 0

func increment(ctx *wasm.HostFunctionCallContext) int32 {
	i += 1
	return i
}

func decrement(ctx *wasm.HostFunctionCallContext) int32 {
	i -= 1
	return i
}

func main() {
	var (
		start time.Time
		d     time.Duration
	)

	buf, err := os.ReadFile("table_export.wasm")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read wasm file: %v", err)
		os.Exit(1)
	}

	modExport, err := wasm.DecodeModule(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode module: %v", err)
		os.Exit(1)
	}

	storeExport := wasm.NewStore(wazeroir.NewEngine())

	err = storeExport.AddHostFunction("go", "increment", reflect.ValueOf(increment))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	err = storeExport.AddHostFunction("go", "decrement", reflect.ValueOf(decrement))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	err = storeExport.Instantiate(modExport, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	buf, err = os.ReadFile("table_test.wasm")
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
	store.ModuleInstances["go"] = &wasm.ModuleInstance{Exports: map[string]*wasm.ExportInstance{}}
	store.ModuleInstances["go"].Exports["tbl"] = storeExport.ModuleInstances[""].Exports["tbl"]
	store.ModuleInstances["go"].Exports["increment"] = storeExport.ModuleInstances["go"].Exports["increment"]
	store.ModuleInstances["go"].Exports["decrement"] = storeExport.ModuleInstances["go"].Exports["decrement"]
	store.ModuleInstances["go"].Exports["wasm_increment"] = storeExport.ModuleInstances[""].Exports["increment"]
	store.ModuleInstances["go"].Exports["wasm_decrement"] = storeExport.ModuleInstances[""].Exports["decrement"]

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	start = time.Now()
	_, _, err = store.CallFunction("", "go_table_test")
	d = time.Now().Sub(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
	fmt.Printf("go_table_test time=%d\n", d.Milliseconds())

	start = time.Now()
	_, _, err = store.CallFunction("", "go_import_test")
	d = time.Now().Sub(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
	fmt.Printf("go_import_test time=%d\n", d.Milliseconds())

	start = time.Now()
	_, _, err = store.CallFunction("", "wasm_table_test")
	d = time.Now().Sub(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
	fmt.Printf("wasm_table_test time=%d\n", d.Milliseconds())

	start = time.Now()
	_, _, err = store.CallFunction("", "wasm_import_test")
	d = time.Now().Sub(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
	fmt.Printf("wasm_import_test time=%d\n", d.Milliseconds())
}
