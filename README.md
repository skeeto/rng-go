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
	"github.com/skeeto/rng-go"
	"math/rand"
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

    # go test -bench=.
    goos: linux
    goarch: amd64
    pkg: github.com/skeeto/rng-go
    BenchmarkLcg128-8                  	1000000000	         2.54 ns/op
    BenchmarkLcg128Interface-8         	300000000	         4.16 ns/op
    BenchmarkSplitMix64-8              	2000000000	         1.51 ns/op
    BenchmarkSplitMix64Interface-8     	300000000	         4.32 ns/op
    BenchmarkXoshiro256ss-8            	500000000	         3.68 ns/op
    BenchmarkXoshiro256ssInterface-8   	200000000	         6.06 ns/op
    BenchmarkPcg32-8                   	300000000	         4.46 ns/op
    BenchmarkPcg32Interface-8          	200000000	         7.12 ns/op
    BenchmarkBaseline-8                	500000000	         3.69 ns/op
    PASS
    ok  	github.com/skeeto/rng-go	19.590s

The big takeaway here: **Interface calls are expensive!** If possible,
use SplitMix64, and do not call it through an interface since that cuts
its performance by 65%. If you must call through an interface, the
built-in PRNG is the fastest. Use it if the large state doesn't matter.
