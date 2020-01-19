# Alternative Go PRNGs

Package `rng` provides several more efficient PRNGs sources for use with
`math/rand.Rand`. Each PRNG implements the `math/rand.Source64`
interface.

What makes these PRNGs more efficient? They have tiny states — 32 bytes
for the largest — compared to gc's default source, which has a ~5kB
state. Two of the generators run faster, too. See the benchmarks below.

What PRNGs are included?

* A ["minimal standard" 128-bit linear congruential generator (LCG)][lcg128]
* [SplitMix64][sm64]
* [xoshiro256\*\*][xo]
* [32-bit permuted congruential generator (PCG)][pcg32]

[lcg128]: http://www.pcg-random.org/posts/does-it-beat-the-minimal-standard.html
[sm64]: http://xoshiro.di.unimi.it/splitmix64.c
[xo]: http://xoshiro.di.unimi.it/xoshiro256starstar.c
[pcg32]: http://www.pcg-random.org/download.html

Requires Go 1.12.

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
    BenchmarkLcg128-8                  	407349361	         2.76 ns/op
    BenchmarkLcg128Interface-8         	312309956	         3.85 ns/op
    BenchmarkSplitMix64-8              	888909228	         1.40 ns/op
    BenchmarkSplitMix64Interface-8     	348746890	         3.47 ns/op
    BenchmarkXoshiro256ss-8            	345606333	         3.50 ns/op
    BenchmarkXoshiro256ssInterface-8   	215023484	         5.58 ns/op
    BenchmarkPcg32-8                   	347164992	         3.54 ns/op
    BenchmarkPcg32Interface-8          	245145093	         4.90 ns/op
    BenchmarkPcg64-8            	222769274	         5.35 ns/op
    BenchmarkPcg64Interface-8   	170089952	         7.05 ns/op
    BenchmarkBaseline-8                	365977761	         3.29 ns/op
    PASS
    ok  	nullprogram.com/x/rng	14.099s

The big takeaway here: **Interface calls are expensive!** If possible,
use SplitMix64, and do not call it through an interface since that cuts
its performance in half. If you must call through an interface, the
built-in PRNG is the fastest, though has the worst quality and a large
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

Tests were run with a zero seed, dieharder 3.31.1, TestU01 1.2.3, and
PractRand 0.95. PractRand was stopped after 8TB of input.
