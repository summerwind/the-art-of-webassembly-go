package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/fatih/color"
	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/wazeroir"
)

func main() {
	var (
		objBaseAddr = 0
		objCount    = 32
		objStride   = 16

		xOffset         = 0
		yOffset         = 4
		radiusOffset    = 8
		collisionOffset = 12

		objI32BaseIndex = objBaseAddr / 4
		objI32Stride    = objStride / 4

		xOffsetI32         = xOffset / 4
		yOffsetI32         = yOffset / 4
		radiusOffsetI32    = radiusOffset / 4
		collisionOffsetI32 = collisionOffset / 4
	)

	buf, err := os.ReadFile("data_structures.wasm")
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

	err = store.AddGlobal("env", "obj_base_addr", uint64(objBaseAddr), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "obj_count", uint64(objCount), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "obj_stride", uint64(objStride), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "x_offset", uint64(xOffset), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "y_offset", uint64(yOffset), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "radius_offset", uint64(radiusOffset), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	err = store.AddGlobal("env", "collision_offset", uint64(collisionOffset), wasm.ValueTypeI32, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add global: %v", err)
		os.Exit(1)
	}

	mem := store.ModuleInstances["env"].Exports["mem"].Memory
	for i := 0; i < objCount; i++ {
		index := objI32Stride*i + objI32BaseIndex
		x := math.Floor(rand.Float64() * 100)
		y := math.Floor(rand.Float64() * 100)
		r := math.Floor(rand.Float64() * 10)
		binary.LittleEndian.PutUint32(mem.Buffer[(index+xOffsetI32)*4:], uint32(x))
		binary.LittleEndian.PutUint32(mem.Buffer[(index+yOffsetI32)*4:], uint32(y))
		binary.LittleEndian.PutUint32(mem.Buffer[(index+radiusOffsetI32)*4:], uint32(r))
	}

	err = store.Instantiate(mod, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate: %v", err)
		os.Exit(1)
	}

	for i := 0; i < objCount; i++ {
		index := objI32Stride*i + objI32BaseIndex
		x := binary.LittleEndian.Uint32(mem.Buffer[(index+xOffsetI32)*4:])
		y := binary.LittleEndian.Uint32(mem.Buffer[(index+yOffsetI32)*4:])
		r := binary.LittleEndian.Uint32(mem.Buffer[(index+radiusOffsetI32)*4:])
		c := binary.LittleEndian.Uint32(mem.Buffer[(index+collisionOffsetI32)*4:])
		if c == 1 {
			color.Red("obj[%02d] x=%2d y=%2d r=%2d collision=%d\n", i, x, y, r, c)
		} else {
			color.Green("obj[%02d] x=%2d y=%2d r=%2d collision=%d\n", i, x, y, r, c)
		}
	}
}
