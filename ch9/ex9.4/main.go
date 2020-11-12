package main

import (
	"flag"
	"fmt"
	"time"
)

var n = flag.Int64("n", 1, "number of pipline stages")

func main() {
	flag.Parse()
	createPipelines(*n)
}

func createPipelines(n int64) {
	first := make(chan int)
	in := first
	var out chan int
	for i := int64(0); i < n; i++ {
		out = createPipeline(in)
		in = out
	}

	start := time.Now()

	first <- 0
	close(first)
	last := <-out

	elapsed := time.Since(start)
	fmt.Printf("n=%d, last=%d, elapsed=%v, ns/stage=%v\n", n, last, elapsed, elapsed.Nanoseconds()/n)
}

func createPipeline(in <-chan int) (out chan int) {
	out = make(chan int)
	go func() {
		n := <-in
		out <- n + 1
		close(out)
	}()
	return
}

/*
$ go run main.go -n 1048576
n=1048576, last=1048576, elapsed=430.479379ms, ns/stage=410
$ go run main.go -n 2097152
n=2097152, last=2097152, elapsed=839.432374ms, ns/stage=400
$ go run main.go -n 4194304
n=4194304, last=4194304, elapsed=3.043621663s, ns/stage=725
$ go run main.go -n 8388608
n=8388608, last=8388608, elapsed=37.609537567s, ns/stage=4483
$ go run main.go -n 16777216
n=16777216, last=16777216, elapsed=1m28.630355116s, ns/stage=5282
*/
