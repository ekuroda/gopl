package popcount

import (
	"testing"
)

func bench(b *testing.B, f func(uint64) int) {
	for i := 0; i < b.N; i++ {
		f(uint64(i))
	}

}

func BenchmarkSingleExp(b *testing.B) {
	bench(b, PopCount)
}

func BenchmarkLoop(b *testing.B) {
	bench(b, PopCountByLoop)
}
