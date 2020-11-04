package popcount

import (
	"fmt"
	"sync"
	"testing"
)

func TestPopCount(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			ch <- PopCount(0x1234567890ABCDEF)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var pcs []int
	for c := range ch {
		pcs = append(pcs, c)
	}

	if got, want := fmt.Sprintf("%v", pcs), fmt.Sprintf("%v", []int{32, 32, 32, 32, 32}); got != want {
		t.Errorf("pcs = %v, want %v", got, want)
	}
}
