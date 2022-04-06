package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
	"time"
)

// funcPerformWasm was compiled from testdata/func_perform.wat
//go:embed testdata/func_perform.wasm
var funcPerformWasm []byte

var i uint32 = 0

func externalCall() uint32 {
	i += 1
	return i
}

func main() {
	r := wazero.NewRuntime()

	goM, err := r.NewModuleBuilder("go").
		ExportFunction("external_call", externalCall).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer goM.Close()

	module, err := r.InstantiateModuleFromCode(funcPerformWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	for _, fn := range []string{"wasm_call", "go_call"} {
		start := time.Now()
		_, err = module.ExportedFunction(fn).Call(nil)
		if err != nil {
			log.Fatal(err)
		}
		d := time.Since(start)

		fmt.Printf("%s time=%d\n", fn, d.Milliseconds())
	}
}
