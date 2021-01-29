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
* [RomuDuo][romu]

Pcg64x is the fastest generator that passes all of the tests.

[lcg128]: http://www.pcg-random.org/posts/does-it-beat-the-minimal-standard.html
[msws]: https://pthree.org/2018/07/30/middle-square-weyl-sequence-prng/
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
    BenchmarkLcg128-8                  	429435103	         2.72 ns/op
    BenchmarkLcg128Interface-8         	305827086	         3.94 ns/op
    BenchmarkSplitMix64-8              	869916241	         1.37 ns/op
    BenchmarkSplitMix64Interface-8     	349949367	         3.40 ns/op
    BenchmarkXoshiro256ss-8            	345413484	         3.49 ns/op
    BenchmarkXoshiro256ssInterface-8   	219676675	         5.58 ns/op
    BenchmarkPcg32-8                   	338895014	         3.54 ns/op
    BenchmarkPcg32Interface-8          	249454053	         4.90 ns/op
    BenchmarkPcg64-8                   	223279998	         5.40 ns/op
    BenchmarkPcg64Interface-8          	169870239	         7.06 ns/op
    BenchmarkPcg64x-8                  	590234182	         2.08 ns/op
    BenchmarkPcg64xInterface-8         	285790071	         4.18 ns/op
    BenchmarkMsws64-8                  	347618887	         3.43 ns/op
    BenchmarkMsws64Interface-8         	245472339	         4.89 ns/op
    BenchmarkRomuDuo-8                 	581715128	         2.03 ns/op
    BenchmarkRomuDuoInterface-8        	331822014	         3.67 ns/op
    BenchmarkBaseline-8                	298460190	         4.05 ns/op

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
| Pcg64x         | PASS      | PASS     | > 8TB     |
| Msws64         | PASS      | PASS     | > 8TB     |
| RomuDuo        | PASS      | PASS     | > 8TB     |

Tests were run with a zero seed, dieharder 3.31.1, TestU01 1.2.3, and
PractRand 0.95. PractRand was stopped after 8TB of input.
