package main

import (
	"encoding/binary"
	"math/rand"
	"os"

	"nullprogram.com/x/rng"
)

func main() {
	// For the faster generators, a closure is faster. For the others,
	// an interface is faster. (gc 1.13)
	var gen func() uint64
	switch os.Args[len(os.Args)-1] {
	case "lcg128":
		gen = new(rng.Lcg128).Uint64
	case "splitmix64":
		gen = new(rng.SplitMix64).Uint64
	case "xoshiro256ss":
		r := new(rng.Xoshiro256ss)
		r.Seed(0)
		gen = r.Uint64
	case "pcg32":
		gen = new(rng.Pcg32).Uint64
	case "pcg64":
		gen = new(rng.Pcg64).Uint64
	case "pcg64x":
		gen = new(rng.Pcg64x).Uint64
	case "msws64":
		gen = new(rng.Msws64).Uint64
	case "baseline":
		gen = rand.NewSource(0).(rand.Source64).Uint64
	default:
		os.Exit(1)
	}

	const n = 1 << 12
	var buf [8 * n]byte
	for {
		for i := 0; i < n; i++ {
			binary.LittleEndian.PutUint64(buf[i*8:], gen())
		}
		if _, err := os.Stdout.Write(buf[:]); err != nil {
			break
		}
	}
}
