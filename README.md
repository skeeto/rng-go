# Alternative Go PRNGs

Package `rng` provides several more efficient, higher quality PRNGs
sources for use with `math/rand.Rand`. Each PRNG implements the
`math/rand.Source64` interface.

What makes these PRNGs more efficient? They have tiny states — 32 bytes
for the largest — compared to gc's default source, which has a ~5kB
state. Two of the generators run faster, too. See the benchmarks below.

What generators are included?

* [SplitMix64][sm64]
* [32-bit and 64-bit permuted congruential generator (PCG)][pcg32]
* Custom 64-bit PCG using [xorshift-multiply][pr] permutation (Pcg64x)
* [xoshiro256\*\*][xo]
* A ["minimal standard" 128-bit linear congruential generator (LCG)][lcg128]
* A 64-bit [Middle Square Weyl Sequence][msws]
* RomuDuo and RomuDuoJr of the [Romu family][romu]
* [Middle Multiplicative Fibonacci Generator][mmlfg]
* [64-bit lag-3 multiply-with-carry generator][mwc256xxa64]

SplitMix64 is the fastest generator. Mmlfg is the fastest robust
generator.

[lcg128]: http://www.pcg-random.org/posts/does-it-beat-the-minimal-standard.html
[mmlfg]: https://github.com/skeeto/scratch/tree/master/mmlfg
[msws]: https://pthree.org/2018/07/30/middle-square-weyl-sequence-prng/
[mwc256xxa64]: https://tom-kaitchuck.medium.com/designing-a-new-prng-1c4ffd27124d
[pcg32]: http://www.pcg-random.org/download.html
[pr]: https://nullprogram.com/blog/2018/07/31/
[romu]: https://romu-random.org/
[sm64]: http://xoshiro.di.unimi.it/splitmix64.c
[xo]: http://xoshiro.di.unimi.it/xoshiro256starstar.c

## Example

```go
package main

import (
	"fmt"
	"math/rand"

	"nullprogram.com/x/rng"
)

func main() {
	s := new(rng.Lcg128)
	s.Seed(1)
	r := rand.New(s)
	fmt.Printf("%v\n", r.NormFloat64())
	fmt.Printf("%v\n", r.NormFloat64())
	fmt.Printf("%v\n", r.NormFloat64())
	// Output:
	// -0.5402515557248266
	// 0.00984877400071782
	// -0.40951475107890106
}
```

## Benchmark

The gc implementation of Go doesn't go a great job optimizing these
routines compared to either GCC or Clang, so SplitMix64 performs the
best of the algorithms in this package. The "baseline" is the default
source from `math/rand`, and the "interface" benchmarks call through the
`math/rand.Source64` interface.

    $ go test -bench=.
    goos: linux
    goarch: amd64
    pkg: nullprogram.com/x/rng
    cpu: Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz
    BenchmarkLcg128-8                  	465076317	         2.521 ns/op
    BenchmarkLcg128Interface-8         	384612382	         3.107 ns/op
    BenchmarkSplitMix64-8              	955716288	         1.252 ns/op
    BenchmarkSplitMix64Interface-8     	391783862	         3.126 ns/op
    BenchmarkXoshiro256ss-8            	498003402	         2.409 ns/op
    BenchmarkXoshiro256ssInterface-8   	381962832	         3.127 ns/op
    BenchmarkPcg32-8                   	421974633	         2.840 ns/op
    BenchmarkPcg32Interface-8          	332382873	         3.598 ns/op
    BenchmarkPcg64-8                   	290823469	         4.146 ns/op
    BenchmarkPcg64Interface-8          	232380673	         5.097 ns/op
    BenchmarkPcg64x-8                  	612730354	         1.957 ns/op
    BenchmarkPcg64xInterface-8         	388036062	         3.095 ns/op
    BenchmarkMsws64-8                  	418388416	         2.999 ns/op
    BenchmarkMsws64Interface-8         	326742526	         3.814 ns/op
    BenchmarkRomuDuo-8                 	647598177	         1.847 ns/op
    BenchmarkRomuDuoInterface-8        	419456901	         2.979 ns/op
    BenchmarkRomuDuoJr-8               	692339223	         1.729 ns/op
    BenchmarkRomuDuoJrInterface-8      	392713688	         3.196 ns/op
    BenchmarkMmlfg-8                   	716652903	         1.605 ns/op
    BenchmarkMmlfgInterface-8          	324306124	         3.730 ns/op
    BenchmarkMwc256xxa64-8             	558840489	         2.153 ns/op
    BenchmarkMwc256xxa64Interface-8    	356378312	         3.454 ns/op
    BenchmarkBaseline-8                	443593838	         2.615 ns/op

The big takeaway here: **Interface calls are relatively expensive!** If
possible, use SplitMix64, and do not call it through an interface since
that cuts its performance in half. If you must call through an interface,
the built-in PRNG is the fastest, though has the worst quality and a large
state.

## Statistical Quality

| generator      | dieharder | BigCrush | PractRand |
|----------------|-----------|----------|-----------|
| math/rand      | PASS      | 1 fail   | 256MB     |
| Lcg128         | PASS      | 1 fail   | 128GB     |
| SplitMix64     | PASS      | 1 fail   | > 8TB     |
| Xoroshiro256ss | PASS      | PASS     | > 8TB     |
| Pcg32          | PASS      | 1 fail   | > 8TB     |
| Pcg64          | PASS      | PASS     | > 8TB     |
| Pcg64x         | PASS      | PASS     | > 8TB     |
| Msws64         | PASS      | PASS     | > 8TB     |
| RomuDuo        | PASS      | PASS     | > 8TB     |
| RomuDuoJr      | PASS      | 3 fail   | > 8TB     |
| Mmlfg          | PASS      | PASS     | > 8TB     |

Tests were run with a zero seed, dieharder 3.31.1, TestU01 1.2.3, and
PractRand 0.95. PractRand was stopped after 8TB of input.
