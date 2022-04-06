package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"log"
)

// stringsWasm was compiled from testdata/strings.wat
//go:embed testdata/strings.wasm
var stringsWasm []byte

func nullStr(m api.Module, strPos uint32) {
	strLen := m.Memory().Size() - strPos
	buf, ok := m.Memory().Read(strPos, strLen)
	if !ok {
		log.Fatalf("Memory.Read(%d, %d) out of range", strPos, strLen)
	}
	str := bytes.Split(buf, []byte{0})
	fmt.Println(string(str[0]))
}

func strPosLen(m api.Module, strPos, strLen uint32) {
	buf, ok := m.Memory().Read(strPos, strLen)
	if !ok {
		log.Fatalf("Memory.Read(%d, %d) out of range", strPos, strLen)
	}
	fmt.Println(string(buf))
}

func lenPrefix(m api.Module, strPos uint32) {
	strLen, ok := m.Memory().ReadByte(strPos)
	if !ok {
		log.Fatalf("Memory.ReadByte(%d) out of range", strPos)
	}
	strPos++
	buf, ok := m.Memory().Read(strPos, uint32(strLen))
	if !ok {
		log.Fatalf("Memory.Read(%d, %d) out of range", strPos, strLen)
	}
	fmt.Println(string(buf))
}

func main() {
	r := wazero.NewRuntime()

	env, err := r.NewModuleBuilder("env").
		ExportFunction("null_str", nullStr).
		ExportFunction("str_pos_len", strPosLen).
		ExportFunction("len_prefix", lenPrefix).
		ExportMemoryWithMax("buffer", 1, 1).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	module, err := r.InstantiateModuleFromCode(stringsWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	test := module.ExportedFunction("main")
	_, err = test.Call(nil)
	if err != nil {
		log.Fatal(err)
	}
}
