package rng_test

import (
	"github.com/skeeto/rng-go"
	"math/rand"
	"testing"
)

func TestLcg128(t *testing.T) {
	want := []uint64{
		0x0fc94e3bf4e9ab32, 0x9f4c53132cb5b55a, 0x04f16bbaa6c209fe,
		0x9c0827f89f0f242f, 0x5b5349ddf2ca0286, 0x9a09a2d3e4f52267,
		0xf4e9e997e821367b, 0xd23cf34fc72f4155, 0x56a2d7e343d7f1b5,
		0x73b5f20e34a8238c, 0xae9a39664ecf3934, 0xe6f5736f43e75071,
		0xf10b6472f469fe94, 0xede9c4aaef957022, 0x8b321466f467bfe0,
	}
	r := rng.Lcg128{0, 1}
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("Lcg128.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func BenchmarkLcg128(b *testing.B) {
	var r rng.Lcg128
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func TestSplitMix64(t *testing.T) {
	want := []uint64{
		0xe220a8397b1dcdaf, 0x6e789e6aa1b965f4, 0x06c45d188009454f,
		0xf88bb8a8724c81ec, 0x1b39896a51a8749b, 0x53cb9f0c747ea2ea,
		0x2c829abe1f4532e1, 0xc584133ac916ab3c, 0x3ee5789041c98ac3,
		0xf3b8488c368cb0a6, 0x657eecdd3cb13d09, 0xc2d326e0055bdef6,
		0x8621a03fe0bbdb7b, 0x8e1f7555983aa92f, 0xb54e0f1600cc4d19,
	}
	r := rng.SplitMix64(0)
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("SplitMix64.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func BenchmarkSplitMix64(b *testing.B) {
	var r rng.SplitMix64
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func TestXoshiro256ss(t *testing.T) {
	// The bad initial output demonstrates why it's important to see
	// this one very carefully. Fortunately it doesn't matter for this
	// test.
	wantInitial := []uint64{
		0x0000000000002d00, 0x0000000000000000, 0x000000005a007080,
		0x10e0000000009d80, 0x10e0b61ce1009d80, 0x0870021ce143ad00,
		0xe071c3c2e143f089, 0x75a1690ef7a20380, 0x9309685b465c23f9,
		0x284f3cc2e13e3c88, 0xc8d749005a413820, 0x1194b410fef20904,
		0xb54a54470263b28c, 0x959e65495daf641c, 0xe561ccecea17f527,
	}
	wantJump := []uint64{
		0x7f7988f72be9c508, 0x5c874fec44783b77, 0x17bcd9b08580dd16,
		0x9ca7f9375f7dbeb2, 0x24caff1483ddd1fa, 0x82d029c9ad74981c,
		0xbbecb7a079cc3631, 0x73e0b137d9f0e369, 0x2b45ddc72e234c08,
		0x06db8f6ecfdb0688, 0xce4ddcf2458e8f71, 0x6892346243ec2224,
		0x721f3bb7498cd45b, 0x4706ddfc3ac5a535, 0x1833b360cae1f78f,
	}
	wantLongJump := []uint64{
		0x409011b83d3299b0, 0xa48dde13c7845f77, 0xf2853b09ce7f46f7,
		0x684e872d5de653df, 0x34d9cef14360b534, 0x42a55e5c647a97c4,
		0xfc07bbe2a0ff49e3, 0x25b74d3c3e1395a4, 0x66c3b4e434a41253,
		0xeef93c334db407df, 0xcbe33255433c267a, 0x1aeb5a580f8b97f7,
		0xee0b16ebb05cc830, 0x1951fff956477d9e, 0xd586fc5de6068234,
	}

	r := rng.Xoshiro256ss{1, 2, 3, 4}
	for i, w := range wantInitial {
		got := r.Uint64()
		if got != w {
			t.Errorf("Xoshiro256ss.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}

	r.Jump()
	for i, w := range wantJump {
		got := r.Uint64()
		if got != w {
			t.Errorf("Xoshiro256ss.Jump(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}

	r.LongJump()
	for i, w := range wantLongJump {
		got := r.Uint64()
		if got != w {
			t.Errorf("Xoshiro256ss.LongJump(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func BenchmarkXoshiro256ss(b *testing.B) {
	var r rng.Xoshiro256ss
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func TestPcg32(t *testing.T) {
	// Output from official "Minimal C Implementation"
	// seed = 0, inc = 0x14057b7ef767814f
	want := []uint32{
		0x00000000, 0x602bf3fd, 0xe823a24e, 0x7a7ecbd9, 0x89fd6c06,
		0xae646aa8, 0xcd3cf945, 0x6204b303, 0x198c8585, 0x49fce611,
		0xd1e9297a, 0x142d9440, 0xee75f56b, 0x473a9117, 0xe3a45903,
		0xbce807a1, 0xe54e5f4d, 0x497d6c51, 0x61829166, 0xa740474b,
		0x031912a8, 0x9de3defa, 0xd266dbf1, 0x0f38bebb, 0xec3c4f65,
	}
	r := rng.Pcg32(0)
	for i, w := range want {
		got := r.Uint32()
		if got != w {
			t.Errorf("Pcg32.Uint32(%d), got %#08x, want %#08x",
				i, got, w)
		}
	}
}

func BenchmarkPcg32(b *testing.B) {
	var r rng.Pcg32
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkBaseline(b *testing.B) {
	// This test isn't entirely fair since it's being done through an
	// interface, but since the concrete implementation isn't exported
	// this is the best we can do!
	r := rand.NewSource(int64(b.N)).(rand.Source64)
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}
