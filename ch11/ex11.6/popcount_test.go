package popcount

import (
	"gopl/ch2/ex2.5/popcount"
	"math"
	"math/rand"
	"testing"
	"time"
)

func bench(b *testing.B, min, max int64, f func(uint64) int) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	rndmax := max - min
	for i := 0; i < b.N; i++ {
		f(uint64(min + rng.Int63n(rndmax)))
	}
}

func BenchmarkSingleExp2(b *testing.B) {
	bench(b, 10, 100, popcount.PopCount)
}

func BenchmarkBitShift2(b *testing.B) {
	bench(b, 10, 100, popcount.PopCountByBitShift)
}

func BenchmarkUnsetBit2(b *testing.B) {
	bench(b, 10, 100, popcount.PopCountByUnsetBit)
}

func BenchmarkSingleExp4(b *testing.B) {
	bench(b, 1000, 10000, popcount.PopCount)
}

func BenchmarkBitShift4(b *testing.B) {
	bench(b, 1000, 10000, popcount.PopCountByBitShift)
}

func BenchmarkUnsetBit4(b *testing.B) {
	bench(b, 1000, 10000, popcount.PopCountByUnsetBit)
}

func BenchmarkSingleExp6(b *testing.B) {
	min := int64(math.Pow10(6))
	bench(b, min, min*10, popcount.PopCount)
}

func BenchmarkBitShift6(b *testing.B) {
	min := int64(math.Pow10(6))
	bench(b, min, min*10, popcount.PopCountByBitShift)
}

func BenchmarkUnsetBit6(b *testing.B) {
	min := int64(math.Pow10(6))
	bench(b, min, min*10, popcount.PopCountByUnsetBit)
}

func BenchmarkSingleExp8(b *testing.B) {
	min := int64(math.Pow10(8))
	bench(b, min, min*10, popcount.PopCount)
}

func BenchmarkBitShift8(b *testing.B) {
	min := int64(math.Pow10(8))
	bench(b, min, min*10, popcount.PopCountByBitShift)
}

func BenchmarkUnsetBit8(b *testing.B) {
	min := int64(math.Pow10(8))
	bench(b, min, min*10, popcount.PopCountByUnsetBit)
}

func BenchmarkSingleExp10(b *testing.B) {
	min := int64(math.Pow10(10))
	bench(b, min, min*10, popcount.PopCount)
}

func BenchmarkBitShift10(b *testing.B) {
	min := int64(math.Pow10(10))
	bench(b, min, min*10, popcount.PopCountByBitShift)
}

func BenchmarkUnsetBit10(b *testing.B) {
	min := int64(math.Pow10(10))
	bench(b, min, min*10, popcount.PopCountByUnsetBit)
}
