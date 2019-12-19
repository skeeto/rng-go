package main

import (
	"bufio"
	"encoding/binary"
	"math/rand"
	"os"

	"github.com/skeeto/rng-go"
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
	case "baseline":
		gen = rand.NewSource(0).(rand.Source64).Uint64
	default:
		os.Exit(1)
	}

	out := bufio.NewWriter(os.Stdout)
	var buf [8]byte
	for {
		binary.LittleEndian.PutUint64(buf[:], gen())
		if _, err := out.Write(buf[:]); err != nil {
			break
		}
	}
}
