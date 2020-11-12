package main

import (
	"fmt"
	"time"
)

func main() {
	first := make(chan int)
	second := make(chan int)
	firstResult := make(chan int)
	secondResult := make(chan int)
	done := make(chan struct{})
	timer := time.NewTimer(1 * time.Second)

	go func() {
		var i int
	L:
		for {
			select {
			case <-done:
				close(second)
				<-first
				firstResult <- i
				break L
			case v, ok := <-first:
				if ok {
					//fmt.Printf("f:%d ", v)
					i++
					second <- v + 1
				}
			default:
			}
		}
	}()

	go func() {
		var i int
	L:
		for {
			select {
			case <-done:
				close(first)
				<-second
				secondResult <- i
				break L
			case v, ok := <-second:
				if ok {
					//fmt.Printf("s:%d ", v)
					i++
					first <- v + 1
				}
			default:
			}
		}
	}()

	first <- 0
	<-timer.C
	close(done)

	c1 := <-firstResult
	c2 := <-secondResult
	fmt.Printf("c1=%d, c2=%d ns/(c1+c2)=%d\n", c1, c2, time.Second.Nanoseconds()/int64(c1+c2))
}
