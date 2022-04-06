package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
)

// storeDataWasm was compiled from testdata/store_data.wat
//go:embed testdata/store_data.wasm
var storeDataWasm []byte

const (
	dataAddr     = 32
	dataCount    = 16
	dataI32Index = dataAddr / 4
)

func main() {
	r := wazero.NewRuntime()

	env, err := r.NewModuleBuilder("env").
		ExportMemoryWithMax("mem", 1, 1).
		ExportGlobalI32("data_addr", dataAddr).
		ExportGlobalI32("data_count", dataCount).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	module, err := r.InstantiateModuleFromCode(storeDataWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	mem := env.ExportedMemory("mem")
	for i := uint32(0); i < dataI32Index+dataCount+4; i++ {
		data, ok := mem.ReadUint32Le(i * 4)
		if !ok {
			log.Fatalf("Memory.ReadUint32Le(%d) out of range", i*4)
		} else if data != 0 {
			fmt.Printf("data[%d]=%v\n", i, data)
		} else {
			fmt.Printf("data[%d]=%v\n", i, data)
		}
	}
}
