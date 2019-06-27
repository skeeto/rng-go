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
routines compared to either GCC or Clang, so SplitMix64 currently
performs the best:

    $ go test -bench=.
    goos: linux
    goarch: amd64
    pkg: github.com/skeeto/rng-go
    BenchmarkLcg128-8               1000000000               2.14 ns/op
    BenchmarkSplitMix64-8           2000000000               1.36 ns/op
    BenchmarkXoshiro256ss-8         500000000                3.48 ns/op
    BenchmarkPcg32-8                300000000                4.21 ns/op
    BenchmarkBaseline-8             500000000                3.29 ns/op
    PASS
    ok      github.com/skeeto/rng-go        10.970s

The "baseline" is the default source from `math/rand`.
