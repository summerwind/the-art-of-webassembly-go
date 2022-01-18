package main

import (
	"fmt"
	"os"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

func main() {
	buf, err := os.ReadFile("pointer.wasm")
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

	maxMem := uint32(1)
	err = store.AddMemoryInstance("env", "mem", 1, &maxMem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add memory instance: %v", err)
		os.Exit(1)
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	results, _, err := store.CallFunction("", "get_ptr")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}

	fmt.Printf("pointer_value=%d\n", results[0])
}
