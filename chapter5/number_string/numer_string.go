package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

func printString(ctx *wasm.HostFunctionCallContext, strPos, strLen int32) {
	buf := ctx.Memory.Buffer[strPos : strPos+strLen]
	fmt.Printf(">%s!\n", buf)
}

func main() {
	value, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid args %v: %v", os.Args[1], err)
		os.Exit(1)
	}

	buf, err := os.ReadFile("number_string.wasm")
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

	err = store.AddHostFunction("env", "print_string", reflect.ValueOf(printString))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	maxMem := uint32(1)
	err = store.AddMemoryInstance("env", "buffer", 1, &maxMem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add memory instance: %v", err)
		os.Exit(1)
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	_, _, err = store.CallFunction("", "to_string", uint64(value))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
}
