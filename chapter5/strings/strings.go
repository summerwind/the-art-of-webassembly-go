package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

var maxMem uint32 = 65535

func nullStr(ctx *wasm.HostFunctionCallContext, strPos int32) {
	buf := ctx.Memory.Buffer[strPos : int32(maxMem)-strPos]
	str := strings.Split(string(buf), string(rune(0)))
	fmt.Println(str[0])
}

func strPosLen(ctx *wasm.HostFunctionCallContext, strPos, strLen int32) {
	buf := ctx.Memory.Buffer[strPos : strPos+strLen]
	fmt.Println(string(buf))
}

func lenPrefix(ctx *wasm.HostFunctionCallContext, strPos int32) {
	len := int(ctx.Memory.Buffer[0:1][0])
	buf := ctx.Memory.Buffer[strPos+1 : len]
	fmt.Println(string(buf))
}

func main() {
	buf, err := os.ReadFile("strings.wasm")
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

	err = store.AddHostFunction("env", "null_str", reflect.ValueOf(nullStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	err = store.AddHostFunction("env", "str_pos_len", reflect.ValueOf(strPosLen))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

	err = store.AddHostFunction("env", "len_prefix", reflect.ValueOf(lenPrefix))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add host function: %v", err)
		os.Exit(1)
	}

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

	_, _, err = store.CallFunction("", "main")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to call function: %v", err)
		os.Exit(1)
	}
}
