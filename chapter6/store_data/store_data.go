package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

func main() {
	var (
		data_addr      = 32
		data_count     = 16
		data_i32_index = data_addr / 4
	)

	buf, err := os.ReadFile("store_data.wasm")
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

	err = store.AddGlobal("env", "data_addr", uint64(data_addr), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "data_count", uint64(data_count), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	mem := store.ModuleInstances["env"].Exports["mem"].Memory
	for i := 0; i < data_i32_index+data_count+4; i++ {
		data := binary.LittleEndian.Uint32(mem.Buffer[i*4:])
		if data != 0 {
			fmt.Printf("data[%d]=%v\n", i, data)
		} else {
			fmt.Printf("data[%d]=%v\n", i, data)
		}
	}
}
