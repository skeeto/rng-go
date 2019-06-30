// Package rng provides several more efficient PRNGs sources for use
// with math/rand.Rand. Each PRNG implements the math/rand.Source64
// interface.
package rng

import (
	"math/bits"
	"math/rand"
)

// An Lcg128 is a 128-bit linear congruential generator implementing
// math/rand.Source64. Note: When seeded manually, the least significant
// bit of Lo must be 1.
type Lcg128 struct{ Hi, Lo uint64 }

var _ rand.Source64 = (*Lcg128)(nil)

func (s *Lcg128) Seed(seed int64) {
	s.Hi = uint64(seed)
	s.Lo = 1
}

func (s *Lcg128) Uint64() uint64 {
	const (
		mhi = 0x0fc94e3bf4e9ab32
		mlo = 0x866458cd56f5e605
	)
	carry, lo := bits.Mul64(mlo, s.Lo)
	s.Hi = mhi*s.Lo + s.Hi*mlo + carry
	s.Lo = lo
	return s.Hi
}

func (s *Lcg128) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// A SplitMix64 provides the SplitMix64 algorithm and implements
// math/rand.Source64. May be manually seeded to any value.
type SplitMix64 uint64

var _ rand.Source64 = (*SplitMix64)(nil)

func (s *SplitMix64) Seed(seed int64) {
	*s = SplitMix64(seed)
}

func (s *SplitMix64) Uint64() uint64 {
	*s += 0x9e3779b97f4a7c15
	z := uint64(*s)
	z = (z ^ z>>30) * 0xbf58476d1ce4e5b9
	z = (z ^ z>>27) * 0x94d049bb133111eb
	return z ^ z>>31
}

func (s *SplitMix64) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// A Xoshiro256ss provides the xoshiro256** algorithm and implements
// math/rand.Source64. Must be seeded carefully with good random values,
// so the Seed() method is highly recommended.
type Xoshiro256ss [4]uint64

var _ rand.Source64 = (*Xoshiro256ss)(nil)

func (s *Xoshiro256ss) Seed(seed int64) {
	var m SplitMix64
	m.Seed(seed)
	s[0] = m.Uint64()
	s[1] = m.Uint64()
	s[2] = m.Uint64()
	s[3] = m.Uint64()
}

func (s *Xoshiro256ss) Uint64() uint64 {
	x := s[1] * 5
	r := (x<<7 | x>>57) * 9
	t := s[1] << 17
	s[2] ^= s[0]
	s[3] ^= s[1]
	s[1] ^= s[2]
	s[0] ^= s[3]
	s[2] ^= t
	s[3] = s[3]<<45 | s[3]>>19
	return r
}

func (s *Xoshiro256ss) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// A Pcg32 provides a 32-bit permuted congruential generator that
// implements math/rand.Source64. Can be seeded to any value.
type Pcg32 uint64

var _ rand.Source64 = (*Pcg32)(nil)

func (s *Pcg32) Seed(seed int64) {
	*s = Pcg32(seed)
	s.Uint32() // discard first output as it's essentially just the seed
}

// Uint32 returns a uniformly random 32-bit integer.
func (s *Pcg32) Uint32() uint32 {
	p := uint64(*s)
	*s = Pcg32(p*0x5851f42d4c957f2d + 0x14057b7ef767814f)
	x := uint32((p>>18 ^ p) >> 27)
	r := uint(p >> 59)
	return x>>r | x<<(-r&31)
}

func (s *Pcg32) Uint64() uint64 {
	lo := uint64(s.Uint32())
	hi := uint64(s.Uint32())
	return hi<<32 | lo
}

func (s *Pcg32) Int63() int64 {
	return int64(s.Uint64() >> 1)
}
