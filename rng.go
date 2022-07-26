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
// implementing math/rand.Source64. Can be seeded to any value. Lcg128
// does not pass PractRand and is here mostly for benchmarking.
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
// math/rand.Source64. Can be seeded to any value.
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
	r := bits.RotateLeft64(x, 7) * 9
	t := s[1] << 17
	s[2] ^= s[0]
	s[3] ^= s[1]
	s[1] ^= s[2]
	s[0] ^= s[3]
	s[2] ^= t
	s[3] = bits.RotateLeft64(s[3], 45)
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
// implements math/rand.Source64. Can be seeded to any value. Pcg32 does
// not pass Big Crush.
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
	r := int(p >> 59)
	return bits.RotateLeft32(x, -r)
}

func (s *Pcg32) Uint64() uint64 {
	lo := uint64(s.Uint32())
	hi := uint64(s.Uint32())
	return hi<<32 | lo
}

func (s *Pcg32) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// A Pcg64 provides a 64-bit permuted congruential generator that
// implements math/rand.Source64. Can be seeded to any value.
type Pcg64 struct{ Hi, Lo uint64 }

var _ rand.Source64 = (*Pcg64)(nil)

func (s *Pcg64) Seed(seed int64) {
	s.Lo = uint64(seed)
	s.Hi = 0
}

func (s *Pcg64) Uint64() uint64 {
	const (
		mhi = 0x2360ed051fc65da4
		mlo = 0x4385df649fccf645
		ahi = 0x5851f42d4c957f2d
		alo = 0x14057b7ef767814f
	)
	carry, lo := bits.Mul64(mlo, s.Lo)
	hi := mhi*s.Lo + s.Hi*mlo + carry
	lo, carry = bits.Add64(lo, alo, 0)
	hi += ahi + carry
	s.Lo = lo
	s.Hi = hi
	lo, hi = lo^lo>>43^hi<<21, hi^hi>>43
	r := int(hi>>60) + 45
	return lo>>r | hi<<(64-r)
}

func (s *Pcg64) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// A Pcg64x provides a 64-bit permuted congruential generator that
// implements math/rand.Source64. Can be seeded to any value. The
// permutation is done with xorshift-multiply. It's much faster than
// Pcg64 but lacks its prediction resistance.
type Pcg64x struct{ Hi, Lo uint64 }

var _ rand.Source64 = (*Pcg64x)(nil)

func (s *Pcg64x) Seed(seed int64) {
	s.Lo = 0xe1cf322879493bf1
	s.Hi = uint64(seed)
}

func (s *Pcg64x) Uint64() uint64 {
	const m = 0xb47d5ba190fb0fa5
	var c uint64
	c, s.Lo = bits.Mul64(s.Lo, m)
	s.Hi = s.Hi*m + c
	s.Lo, c = bits.Add64(s.Lo, 1, 0)
	s.Hi += c
	r := s.Hi
	r ^= r >> 32
	r *= m
	return r
}

func (s *Pcg64x) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// A Msws64 is the Middle Square Weyl Sequence algorithm. It implements
// math/rand.Source64 and may be seeded to any value.
type Msws64 [4]uint64

var _ rand.Source64 = (*Msws64)(nil)

func (s *Msws64) Seed(seed int64) {
	v := uint64(seed)
	*s = Msws64{v, v, v, v}
}

func (s *Msws64) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *Msws64) Uint64() uint64 {
	var xl, xh, wl, wh, c uint64
	c, xl = bits.Mul64(s[0], s[0])
	xh = 2*s[0]*s[1] + c
	wl, c = bits.Add64(s[2], 0x8367589d496e8afd, 0)
	wh = s[3] + 0x918fba1eff8e67e1 + c
	xl, c = bits.Add64(xl, wl, 0)
	xh = xh + wh + c
	s[0] = xh
	s[1] = xl
	s[2] = wl
	s[3] = wh
	return xh
}

// A RomuDuo is a chaotic generator that combines the linear operation
// of multiplication with the nonlinear operation of rotation. Must be
// seeded carefully with good random values, so the Seed() method is
// highly recommended.
type RomuDuo struct{ x, y uint64 }

var _ rand.Source64 = (*RomuDuo)(nil)

func (s *RomuDuo) Seed(seed int64) {
	var m SplitMix64
	m.Seed(seed)
	s.x = m.Uint64()
	s.y = m.Uint64()
}

func (s *RomuDuo) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *RomuDuo) Uint64() uint64 {
	x := s.x
	s.x = 0xd3833e804f4c574b * s.y
	s.y = bits.RotateLeft64(s.y, 36) + bits.RotateLeft64(s.y, 15) - x
	return x
}

// A RomuDuoJr is a chaotic generator that combines the linear operation
// of multiplication with the nonlinear operation of rotation. It should
// be slighter faster than RomuDuoJr at the cost of reduced capacity.
// Must be seeded carefully with good random values, so the Seed()
// method is highly recommended.
type RomuDuoJr struct{ x, y uint64 }

var _ rand.Source64 = (*RomuDuoJr)(nil)

func (s *RomuDuoJr) Seed(seed int64) {
	var m SplitMix64
	m.Seed(seed)
	s.x = m.Uint64()
	s.y = m.Uint64()
}

func (s *RomuDuoJr) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *RomuDuoJr) Uint64() uint64 {
	x := s.x
	s.x = 0xd3833e804f4c574b * s.y
	s.y -= x
	s.y = bits.RotateLeft64(s.y, 27)
	return x
}

// A Mmlfg is a Middle Multiplicative Lagged Fibonacci generator. The
// output is the middle 64 bits of a 128-bit product. A larger state is
// required to pass statistical tests. Must be seeded carefully with
// good random values, and all state elements must be odd, so the Seed()
// method is highly recommended.
type Mmlfg struct {
	s    [15]uint64
	i, j int32
}

var _ rand.Source = (*Mmlfg)(nil)

func (m *Mmlfg) Seed(seed int64) {
	s := uint64(seed)
	for i := 0; i < 15; i++ {
		s = s*0x3243f6a8885a308d + 1111111111111111111
		m.s[i] = s ^ s>>31 | 1
	}
	m.i = 14
	m.j = 12
}

func (m *Mmlfg) Int63() int64 {
	return int64(m.Uint64() >> 1)
}

func (m *Mmlfg) Uint64() uint64 {
	hi, lo := bits.Mul64(m.s[m.i], m.s[m.j])
	m.s[m.i] = lo
	m.i--
	if m.i < 0 {
		m.i = 14
	}
	m.j--
	if m.j < 0 {
		m.j = 14
	}
	return hi<<32 | lo>>32
}

// A Mwc256xxa64 is a 64-bit lag-3 multiply-with-carry generator. Must
// be seeded carefully with good random values, so the Seed() method is
// highly recommended.
type Mwc256xxa64 [4]uint64

var _ rand.Source64 = (*Mwc256xxa64)(nil)

func (m *Mwc256xxa64) Seed(seed int64) {
	m[0] = uint64(seed)
	m[1] = uint64(seed)
	m[2] = 0xcafef00dd15ea5e5
	m[3] = 0x14057b7ef767814f
	for i := 0; i < 6; i++ {
		m.Uint64()
	}
}

func (m *Mwc256xxa64) Int63() int64 {
	return int64(m.Uint64() >> 1)
}

func (m *Mwc256xxa64) Uint64() uint64 {
	hi, lo := bits.Mul64(0xfeb344657c0af413, m[2])
	r := (m[2] ^ m[1]) + (m[0] ^ hi)
	t, c := bits.Add64(m[3], lo, 0)
	m[2] = m[1]
	m[1] = m[0]
	m[0] = t
	m[3] = hi + c
	return r
}

// An Sfc64 is a 64-bit "small, fast, chaotic" generator. Must be seeded
// carefully with good random values, so the Seed() method is highly
// recommended.
type Sfc64 [4]uint64

var _ rand.Source64 = (*Sfc64)(nil)

func (s *Sfc64) Seed(seed int64) {
	var m SplitMix64
	m.Seed(seed)
	s[0] = m.Uint64()
	s[1] = m.Uint64()
	s[2] = m.Uint64()
	s[3] = m.Uint64()
}

func (s *Sfc64) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *Sfc64) Uint64() uint64 {
	r := s[0] + s[1] + s[3]
	s[3]++
	s[0] = (s[1] >> 11) ^ s[1]
	s[1] = (s[2] << 3) + s[2]
	s[2] = r + (s[2]<<24 | s[2]>>40)
	return r
}
