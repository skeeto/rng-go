package rng_test

import (
	"math/rand"
	"testing"

	"nullprogram.com/x/rng"
)

func TestLcg128(t *testing.T) {
	want := []uint64{
		0x2d99787926d46932, 0x579d64f7b4780f53, 0xc716c8bffcc60271,
		0xfc763fac42f18290, 0xeba26e07402a33f4, 0xe2c6dd9f0e06fc35,
		0x779d001d1e3bf290, 0xfaa9b1ae526c3070, 0x235d2825e14c0f15,
		0x19c3b1bfec64fa79, 0x9ae3d4f0ade39da9, 0x597c849c597c0624,
		0x8be750f54de1d4c4, 0x58d34b21dc53606e, 0x5f78dea7e0db0986,
	}
	r := rng.Lcg128{0, 0}
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

func BenchmarkLcg128Interface(b *testing.B) {
	r := rand.New(new(rng.Lcg128))
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

func BenchmarkSplitMix64Interface(b *testing.B) {
	r := rand.New(new(rng.SplitMix64))
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

func BenchmarkXoshiro256ssInterface(b *testing.B) {
	r := rand.New(new(rng.Xoshiro256ss))
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

func TestPcg64(t *testing.T) {
	want := []uint64{
		0x7d0b53279d44e46b, 0xdf64cca5bee7ba96, 0xae79c64d0ddeef66,
		0x4b119973e1d3c11b, 0xb1ba34b04ea86b85, 0x8bf94307ca1db73b,
		0xb0a32852cf9b69f3, 0xb5e47e28e9159092, 0xf0f56e7cbcd8a441,
		0x71fa9ce37e8d1e62, 0xf3381e9ed062742c, 0xef6e8f4c998a2723,
		0xbdb435240c1b06d3, 0x2e3fe044f1324b00, 0xa70ce29c1bbad6c4,
		0x970f930df1818e2d, 0xcfc29743712bb28d, 0xbb1e68d716ab35a5,
		0x14a5503d53c2c201, 0xba18370ae44e4980, 0xcea50b483a4d2235,
		0x68b655cd5065fc54, 0xfd8a7c265421aa21, 0x2c77dc12c533bb60,
		0x27e3c5efff604f45, 0x0daa271bc7685814, 0xe27337a6cea866f4,
		0x89236f2e868409c9, 0x3a8c12785dd98c9c, 0xf95f97baa171fcd8,
	}
	var r rng.Pcg64
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("Pcg64.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func TestPcg64x(t *testing.T) {
	want := []uint64{
		0x3df9dcd05ccda305, 0x16f6db58022bacc1, 0xa63b6362a3cd40f7,
		0xca13dfd56ea8ef4e, 0x5853275645fe74b5, 0x4ba45d105546008c,
		0xfd9e9d8f90a90c56, 0xbf25345a63d46bd3, 0xc22456f591686cf3,
		0x0b008f2ede6e0d91, 0x501d1cc4c0ba52a7, 0xa0d7998490c022c6,
		0xcf28587e3252172c, 0x9cbeb645cba65c2a, 0x8ca3f7949a2607de,
		0xc202ce927768aa15, 0x90e0fde8ceb6ca7e, 0xaf0dc1093a3e0e44,
		0x34cbad7867c2755c, 0x1af455737add206b, 0xff312e21f55a694a,
		0xa241a4e6f07bba19, 0x3df6845400364da4, 0xb0981cb89a3253f5,
		0xff9ba1711ec5dc8f, 0x306a3658b01393e4, 0xe5aaf3e5692b4650,
		0x91741b16cdba1f1a, 0x70ae0a10095443d4, 0x729d4e374832f0d6,
	}
	var r rng.Pcg64x
	r.Seed(0)
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("Pcg64x.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func TestMsws64(t *testing.T) {
	want := []uint64{
		0x918fba1eff8e67e1, 0xfa29d22ad162985d, 0x545331c2f0cb035c,
		0xd03ec88f2d382bfe, 0xdb3c827bbd00c8c4, 0x22a259673ee1be12,
		0x298520849901be70, 0x3a846be2945f61ea, 0xa56420383c942a5a,
		0xd4aeb5a64faf6dfb, 0x233872349fe640cf, 0x493581b803e06b51,
		0xd99ea19ace03315c, 0x9c72ca6deac41460, 0x1f3699ad4987e7dd,
		0x1405f01787af8b14, 0xd6fdd2c7f34e9f0b, 0x7604bd520718d67f,
		0xc15fc1b374ce737d, 0x9038fe3a45fc4e36, 0xb2eecd0c7da8f67f,
		0x3ef2f6aec20f18e6, 0xb460c5be24073552, 0xfc55b398b9413510,
		0xb2ab678dd49f9078, 0x1d6cf3d7ab040855, 0xbca32286c0f18121,
		0x374d5a547dcb0bdb, 0xbccadd0aa6bd2d9c, 0x9a4edce81329e4ba,
	}
	var r rng.Msws64
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("Msws64.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func TestRomuDuo(t *testing.T) {
	// Output from official reference implementation
	want := []uint64{
		0xe220a8397b1dcdaf, 0x55fcf1b3f366ca7c, 0x9eb834e3c190311d,
		0xab9d6b6df8e3c1be, 0x90b4c67e093a6504, 0xe10add1516363e61,
		0xfffc0e0cda58548b, 0x6029fb83b007ecca, 0x3df83004acca9c4d,
		0x5a17629d93412209, 0x3d7b6e817c1a1c74, 0x53d4d21a8fb0ace7,
		0x4482d635c3a1c5c8, 0x38d7151bf90119b3, 0xb704bcee1abc049a,
		0xf9b01471cc2f51fe, 0x56d8a501d769675e, 0x78ae74325e5bf737,
		0x2d1df34e617c96c0, 0xd5a61bd9c826e5d8, 0x4578ebd26b216724,
		0x062b8632b54a0780, 0xed1fd50dd91e4a33, 0xcfcf5a348f47a172,
		0x02ddda770f9f2dc2, 0x2cf440df9e5d064d, 0xf3d13d9148b0b382,
		0xf81daac0b171e188, 0x14c08d919d08ff9e, 0x654891bee7b6b17e,
	}
	var r rng.RomuDuo
	r.Seed(0)
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("RomuDuo.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func TestRomuDuoJr(t *testing.T) {
	// Output from official reference implementation
	want := []uint64{
		0xe220a8397b1dcdaf, 0x55fcf1b3f366ca7c, 0xb53a06f1179f4fdb,
		0x7f84f708e631f6c8, 0xd11049a7010d66b7, 0x05dbe41727bc9981,
		0x67eeb580b5b35570, 0x64bb861864910769, 0xaf4d8e9876f962d8,
		0x8ef0c65faf9fa5d1, 0xa7c45111ce04ee51, 0xb23394f37a2b2e16,
		0x8a1c4e76add7024a, 0x3888050ae21b4790, 0xded5953f7b0f982e,
		0x6d32a54e205eed8e, 0x69b13d9dc6c3e587, 0xe03e20f8711330fb,
		0x7c59cfc897d11c1f, 0x941cfc090a173196, 0x12b7f1eb1e2e1dea,
		0x2c47a8278953b216, 0xbc9bf8d5dcf3fb03, 0x159d18fed204d64f,
		0x7e8fdbf00694d1c4, 0xbc7237d5ffa1603e, 0x5e32c67921c91b04,
		0xd374c900adb2d55f, 0x08951eef6af92fae, 0x0169c0c46f15b4e1,
	}
	var r rng.RomuDuoJr
	r.Seed(0)
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("RomuDuo.Uint64(%d), got %#016x, want %#016x",
				i, got, w)
		}
	}
}

func TestMmlfg(t *testing.T) {
	// Output from Lua implementation
	want := []uint64{
		0x1573aa52f814bda8, 0x3aeaac28b52676e2, 0x8f1b6491309e5792,
		0x25bca26e169f58cd, 0xee13266f6d5bad81, 0xd688681022995579,
		0xc227f64fffc6967a, 0x3d06e4f91995745f, 0x4077b1108d5150b1,
		0x41deb8bcf496aac3, 0xdef5ecadb01c5527, 0x42be0306aca9476d,
		0xcc40df9abc49fae2, 0xd6fab4fe6f2c8373, 0xad02822ecc846c6d,
		0x602b2201cc7bf7b7, 0xded4343bd0724597, 0xfcbcd8d91b8f65f4,
		0xfc76214430f94e44, 0x4c7fc6e9f4291294, 0xfca3ad5722cee412,
		0xe3383e408585396a, 0xfbafa05b7c2faecf, 0xe684088050284b8c,
		0x8bbb114ed18162a0, 0x0bbde9b2d192d39b, 0xb403be5f2fb967e5,
		0xc60ea291e01fe627, 0x1790ba5d87432edc, 0x598bdded3fe137d9,
		0x0dba6bcb0e9e17ef, 0x748d4dac10754ca0, 0xa212d97e7982de85,
		0x975ea1c76b0f0a7e, 0xad0170d0b44d8673, 0xa3d8fb24e994e7cf,
		0x5ecef8bd9f6e7279, 0xc3a57186c73c6a98, 0x7f3ad93171dfdff9,
	}
	var r rng.Mmlfg
	r.Seed(0)
	for i, w := range want {
		got := r.Uint64()
		if got != w {
			t.Errorf("Mmlfg.Uint64(%d), got %#016x, want %#016x",
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

func BenchmarkPcg32Interface(b *testing.B) {
	r := rand.New(new(rng.Pcg32))
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkPcg64(b *testing.B) {
	var r rng.Pcg64
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkPcg64Interface(b *testing.B) {
	r := rand.New(new(rng.Pcg64))
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkPcg64x(b *testing.B) {
	var r rng.Pcg64x
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkPcg64xInterface(b *testing.B) {
	r := rand.New(new(rng.Pcg64x))
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkMsws64(b *testing.B) {
	var r rng.Msws64
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkMsws64Interface(b *testing.B) {
	r := rand.New(new(rng.Msws64))
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkRomuDuo(b *testing.B) {
	var r rng.RomuDuo
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkRomuDuoInterface(b *testing.B) {
	r := rand.New(new(rng.RomuDuo))
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkRomuDuoJr(b *testing.B) {
	var r rng.RomuDuoJr
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkRomuDuoJrInterface(b *testing.B) {
	r := rand.New(new(rng.RomuDuoJr))
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkMmlfg(b *testing.B) {
	var r rng.Mmlfg
	r.Seed(int64(b.N))
	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkMmlfgInterface(b *testing.B) {
	r := rand.New(new(rng.Mmlfg))
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
