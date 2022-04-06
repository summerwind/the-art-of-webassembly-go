package main

import (
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"log"
	"time"
)

// tableExportWasm was compiled from testdata/table_export.wat
//go:embed testdata/table_export.wasm
var tableExportWasm []byte

// tableTestWasm was compiled from testdata/table_test.wat
//go:embed testdata/table_test.wasm
var tableTestWasm []byte

var i int32 = 0

func increment() int32 {
	i += 1
	return i
}

func decrement() int32 {
	i -= 1
	return i
}

func main() {
	r := wazero.NewRuntime()

	goM, err := r.NewModuleBuilder("go").
		ExportFunction("increment", increment).
		ExportFunction("decrement", decrement).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer goM.Close()

	// Instantiate the tableExportWasm assigning its module name to "wasm"
	tableExport, err := r.InstantiateModuleFromCodeWithConfig(tableExportWasm,
		wazero.NewModuleConfig().WithName("wasm"))
	if err != nil {
		log.Fatal(err)
	}
	defer tableExport.Close()

	// table_test.wat requires all imports to be under the module "go", so we have to replace them.
	module, err := r.InstantiateModuleFromCodeWithConfig(tableTestWasm, wazero.NewModuleConfig().
		WithImport("go", "tbl", "wasm", "tbl").
		WithImport("go", "wasm_increment", "wasm", "increment").
		WithImport("go", "wasm_decrement", "wasm", "decrement"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	for _, fn := range []string{"go_table_test", "go_import_test", "wasm_table_test", "wasm_import_test"} {
		start := time.Now()
		_, err = module.ExportedFunction(fn).Call(nil)
		if err != nil {
			log.Fatal(err)
		}
		d := time.Since(start)

		fmt.Printf("%s time=%d\n", fn, d.Milliseconds())
	}
}
