package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

var startStringIndex uint32 = 100

func prinntString(ctx *wasm.HostFunctionCallContext, strLen uint32) {
	buf := ctx.Memory.Buffer[startStringIndex : startStringIndex+strLen]
	fmt.Printf("%s\n", buf)
}

func main() {
	buf, err := os.ReadFile("helloworld.wasm")
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

	err = store.AddHostFunction("env", "print_string", reflect.ValueOf(prinntString))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	memLimitMax := uint32(1)
	err = store.AddMemoryInstance("env", "buffer", 1, &memLimitMax)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add memory instance: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "start_string", uint64(startStringIndex), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	_, _, err = store.CallFunction("", "helloworld")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
}
