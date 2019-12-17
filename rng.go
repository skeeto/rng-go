// This is free and unencumbered software released into the public domain.

// Package rng provides several more efficient PRNGs sources for use
// with math/rand.Rand. Each PRNG implements the math/rand.Source64
// interface.
package rng

import (
	"math/bits"
	"math/rand"
)

// An Lcg128 is a truncated 128-bit linear congruential generator
// implementing math/rand.Source64. Can be seeded to any value.
type Lcg128 struct{ Hi, Lo uint64 }

var _ rand.Source64 = (*Lcg128)(nil)

func (s *Lcg128) Seed(seed int64) {
	s.Lo = uint64(seed)
	s.Hi = 0
}

func (s *Lcg128) Uint64() uint64 {
	const (
		mhi = 0x2d99787926d46932
		mlo = 0xa4c1f32680f70c55
	)
	carry, lo := bits.Mul64(mlo, s.Lo)
	hi := mhi*s.Lo + s.Hi*mlo + carry
	lo, carry = bits.Add64(lo, mlo, 0)
	hi += mhi + carry
	s.Lo = lo
	s.Hi = hi
	return hi
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
	z ^= z >> 30
	z *= 0xbf58476d1ce4e5b9
	z ^= z >> 27
	z *= 0x94d049bb133111eb
	z ^= z >> 31
	return z
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

var jump = [4]uint64{
	0x180ec6d33cfd0aba, 0xd5a61266f0c9392c,
	0xa9582618e03fc9aa, 0x39abdc4529b1661c,
}

// Jump is equivalent to 2^128 calls to Uint64().
func (s *Xoshiro256ss) Jump() {
	var s0, s1, s2, s3 uint64
	for _, j := range jump {
		for b := uint(0); b < 64; b++ {
			if j&(1<<b) != 0 {
				s0 ^= s[0]
				s1 ^= s[1]
				s2 ^= s[2]
				s3 ^= s[3]
			}
			s.Uint64()
		}
	}
	s[0] = s0
	s[1] = s1
	s[2] = s2
	s[3] = s3
}

var longjump = [4]uint64{
	0x76e15d3efefdcbbf, 0xc5004e441c522fb3,
	0x77710069854ee241, 0x39109bb02acbe635,
}

// LongJump is equivalent to 2^192 calls to Uint64().
func (s *Xoshiro256ss) LongJump() {
	var s0, s1, s2, s3 uint64
	for _, j := range longjump {
		for b := uint(0); b < 64; b++ {
			if j&(1<<b) != 0 {
				s0 ^= s[0]
				s1 ^= s[1]
				s2 ^= s[2]
				s3 ^= s[3]
			}
			s.Uint64()
		}
	}
	s[0] = s0
	s[1] = s1
	s[2] = s2
	s[3] = s3
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
