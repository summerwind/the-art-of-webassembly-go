package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

var i uint32 = 0

func externalCall(ctx *wasm.HostFunctionCallContext) uint32 {
	i += 1
	return i
}

func main() {
	var (
		start time.Time
		d     time.Duration
	)

	buf, err := os.ReadFile("func_perform.wasm")
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

	err = store.AddHostFunction("go", "external_call", reflect.ValueOf(externalCall))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	start = time.Now()
	_, _, err = store.CallFunction("", "wasm_call")
	d = time.Now().Sub(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
	fmt.Printf("wasm_call time=%d\n", d.Milliseconds())

	start = time.Now()
	_, _, err = store.CallFunction("", "go_call")
	d = time.Now().Sub(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
	fmt.Printf("go_call time=%d\n", d.Milliseconds())
}
