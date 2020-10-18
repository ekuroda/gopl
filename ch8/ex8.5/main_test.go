package main

import "testing"

func benchmark(b *testing.B, workerNum int) {
	for i := 0; i < b.N; i++ {
		createImage(workerNum)
	}
}

func Benchmark1(b *testing.B)  { benchmark(b, 1) }
func Benchmark2(b *testing.B)  { benchmark(b, 2) }
func Benchmark4(b *testing.B)  { benchmark(b, 4) }
func Benchmark6(b *testing.B)  { benchmark(b, 6) }
func Benchmark8(b *testing.B)  { benchmark(b, 8) }
func Benchmark10(b *testing.B) { benchmark(b, 10) }
func Benchmark12(b *testing.B) { benchmark(b, 12) }
func Benchmark14(b *testing.B) { benchmark(b, 14) }
func Benchmark16(b *testing.B) { benchmark(b, 16) }
func Benchmark20(b *testing.B) { benchmark(b, 20) }
