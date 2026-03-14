package main

import (
	"chess/internal/game"
	// "runtime"
	"os"
	"runtime/pprof"
)

func main() {
	f, _ := os.Create("cpu.prof")
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	game.Start()
}
