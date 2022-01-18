package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

func log(ctx *wasm.HostFunctionCallContext, n, factorial uint32) {
	fmt.Printf("%d! = %d\n", n, factorial)
}

func main() {
	n, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid args %v: %v", os.Args[1], err)
		os.Exit(1)
	}

	buf, err := os.ReadFile("loop.wasm")
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

	err = store.AddHostFunction("env", "log", reflect.ValueOf(log))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	results, _, err := store.CallFunction("", "loop_test", n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}

	fmt.Printf("result %d! = %d\n", n, results[0])
	if n > 12 {
		fmt.Println("Factorials greater than 12 are too large for a 32-bit integer.")
	}
}
