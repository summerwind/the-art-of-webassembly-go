package main

import (
	"github.com/fatih/color"
	"os"
)

// ExampleMain ensures the following will work:
//
//	go build data_structures.go
//	./data_structures
func ExampleMain() {
	color.NoColor = true
	color.Output = os.Stdout

	main()

	// Output:
	// obj[00] x=60 y=94 r= 6 collision=1
	// obj[01] x=43 y=42 r= 6 collision=0
	// obj[02] x= 6 y=15 r= 0 collision=0
	// obj[03] x=30 y=51 r= 8 collision=1
	// obj[04] x=21 y=38 r= 3 collision=0
	// obj[05] x=46 y=28 r= 2 collision=0
	// obj[06] x=67 y=21 r= 2 collision=1
	// obj[07] x=36 y=57 r= 8 collision=1
	// obj[08] x=29 y=29 r= 7 collision=0
	// obj[09] x=20 y=86 r= 6 collision=0
	// obj[10] x=52 y= 2 r= 1 collision=0
	// obj[11] x=60 y=97 r= 0 collision=1
	// obj[12] x=59 y= 5 r= 6 collision=1
	// obj[13] x=30 y=17 r= 5 collision=1
	// obj[14] x=54 y=27 r= 4 collision=1
	// obj[15] x=53 y=25 r= 2 collision=1
	// obj[16] x=78 y=36 r= 8 collision=1
	// obj[17] x=29 y=89 r= 0 collision=0
	// obj[18] x=97 y= 7 r= 2 collision=1
	// obj[19] x=68 y=24 r= 3 collision=1
	// obj[20] x=93 y=74 r= 8 collision=1
	// obj[21] x=73 y=18 r= 4 collision=0
	// obj[22] x=89 y=68 r= 9 collision=1
	// obj[23] x=92 y= 9 r= 4 collision=1
	// obj[24] x=92 y=95 r= 3 collision=0
	// obj[25] x=69 y=71 r= 5 collision=0
	// obj[26] x=64 y=55 r= 7 collision=0
	// obj[27] x=40 y=13 r= 9 collision=1
	// obj[28] x=89 y=32 r= 7 collision=1
	// obj[29] x=64 y= 8 r= 6 collision=1
	// obj[30] x=62 y=36 r= 2 collision=0
	// obj[31] x=53 y=18 r= 2 collision=0
}
