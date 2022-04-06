package main

import (
	_ "embed"
	"github.com/fatih/color"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"log"
	"math"
	"math/rand"
)

// dataStructuresWasm was compiled from testdata/data_structures.wat
//go:embed testdata/data_structures.wasm
var dataStructuresWasm []byte

const (
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

func main() {
	r := wazero.NewRuntime()

	env, err := r.NewModuleBuilder("env").
		ExportMemoryWithMax("mem", 1, 1).
		ExportGlobalI32("obj_base_addr", objBaseAddr).
		ExportGlobalI32("obj_count", objCount).
		ExportGlobalI32("obj_stride", objStride).
		ExportGlobalI32("x_offset", xOffset).
		ExportGlobalI32("y_offset", yOffset).
		ExportGlobalI32("radius_offset", radiusOffset).
		ExportGlobalI32("collision_offset", collisionOffset).
		Instantiate()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	mem := env.ExportedMemory("mem")
	for i := uint32(0); i < objCount; i++ {
		index := objI32Stride*i + objI32BaseIndex
		x := uint32(math.Floor(rand.Float64() * 100))
		y := uint32(math.Floor(rand.Float64() * 100))
		r := uint32(math.Floor(rand.Float64() * 10))

		xO := (index + xOffsetI32) * uint32(4)
		yO := (index + yOffsetI32) * uint32(4)
		rO := (index + radiusOffsetI32) * uint32(4)

		for _, t := range [][2]uint32{{xO, x}, {yO, y}, {rO, r}} {
			if ok := mem.WriteUint32Le(t[0], t[1]); !ok {
				log.Fatalf("Memory.WriteUint32Le(%d, %d) out of range", t[0], t[1])
			}
		}
	}

	module, err := r.InstantiateModuleFromCode(dataStructuresWasm)
	if err != nil {
		log.Fatal(err)
	}
	defer module.Close()

	for i := uint32(0); i < objCount; i++ {
		index := objI32Stride*i + objI32BaseIndex

		xO := (index + xOffsetI32) * uint32(4)
		yO := (index + yOffsetI32) * uint32(4)
		rO := (index + radiusOffsetI32) * uint32(4)
		cO := (index + collisionOffsetI32) * uint32(4)

		x := readUint32Le(mem, xO)
		y := readUint32Le(mem, yO)
		r := readUint32Le(mem, rO)
		c := readUint32Le(mem, cO)

		if c == 1 {
			color.Red("obj[%02d] x=%2d y=%2d r=%2d collision=%d\n", i, x, y, r, c)
		} else {
			color.Green("obj[%02d] x=%2d y=%2d r=%2d collision=%d\n", i, x, y, r, c)
		}
	}
}

func readUint32Le(mem api.Memory, offset uint32) uint32 {
	v, ok := mem.ReadUint32Le(offset)
	if !ok {
		log.Fatalf("Memory.ReadUint32Le(%d) out of range", offset)
	}
	return v
}
