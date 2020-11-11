package memo

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

type resultData struct {
	key string
	err error
	t   time.Duration
}

func TestGet(t *testing.T) {
	m := New(delay)
	defer m.Close()

	keys := []string{
		"1", "2", "3", "1", "2", "3",
	}
	doneShort := make(chan struct{})
	doneLong := make(chan struct{})
	dones := []chan struct{}{
		doneShort, doneLong, doneShort, doneLong, doneLong, doneShort,
	}

	errorCount := make(map[string]int)
	elapsedTimes := make(map[string][]time.Duration)
	resultChan := make(chan resultData)

	go func() {
		for r := range resultChan {
			if r.err != nil {
				errorCount[r.key]++
			} else {
				val := elapsedTimes[r.key]
				val = append(val, r.t)
				elapsedTimes[r.key] = val
			}
		}
	}()

	var n sync.WaitGroup
	n.Add(1)
	go func() {
		defer n.Done()
		timer := time.NewTimer(500 * time.Millisecond)
		<-timer.C
		log.Printf("done short")
		close(doneShort)
	}()

	n.Add(1)
	go func() {
		defer n.Done()
		timer := time.NewTimer(2000 * time.Millisecond)
		<-timer.C
		log.Printf("done long")
		close(doneLong)
	}()

	for i, key := range keys {
		n.Add(1)
		go func(key string, done chan struct{}) {
			defer n.Done()
			start := time.Now()
			_, err := m.Get(key, done)
			elapsed := time.Since(start)
			if err != nil {
				log.Print(err)
			} else {
				log.Printf("%s, %s", key, elapsed)
			}
			resultChan <- resultData{key, err, elapsed}
		}(key, dones[i])
	}

	timer := time.NewTimer(600 * time.Millisecond)
	<-timer.C

	keys = []string{"3", "3"}
	dones = []chan struct{}{doneLong, doneLong}
	for i, key := range keys {
		n.Add(1)
		go func(key string, done chan struct{}) {
			defer n.Done()
			start := time.Now()
			_, err := m.Get(key, done)
			elapsed := time.Since(start)
			if err != nil {
				log.Print(err)
			} else {
				log.Printf("%s, %s", key, elapsed)
			}
			resultChan <- resultData{key, err, elapsed}
		}(key, dones[i])
	}

	n.Wait()
	close(resultChan)

	log.Printf("%v, %v", elapsedTimes, errorCount)

	wants := []struct {
		key         string
		resultCount int
		errorCount  int
	}{
		{"1", 1, 1},
		{"2", 2, 0},
		{"3", 2, 2},
	}
	for _, w := range wants {
		k := w.key
		v, ok := elapsedTimes[k]
		if !ok {
			t.Errorf("elapsedTimes[%s] not exists", k)
			continue
		}
		if len(v) != w.resultCount {
			t.Errorf("elapsedTimes[%s] has %d items, want %d", k, w.resultCount, len(v))
			continue
		}
		for i, item := range v {
			if item <= time.Second || 2*time.Second <= item {
				t.Errorf("elapsedTimes[%s][%d] = %v, want between 1s and 2s", k, i, item)
			}
		}

		ec := errorCount[k]
		if ec != w.errorCount {
			t.Errorf("errorCount[%s] = %d, want %d", k, ec, w.errorCount)
		}
	}
}

func delay(key string, done <-chan struct{}) (result interface{}, err error) {
	timer := time.NewTimer(time.Second)
	exit := false
	for !exit {
		select {
		case <-timer.C:
			exit = true
			result = true
		case <-done:
			exit = true
			err = fmt.Errorf("canceled")
		}
	}
	return
}

// func TestGet(t *testing.T) {
// 	m := New(httpGetBody)
// 	done := make(chan struct{})
// 	defer m.Close()

// 	var n sync.WaitGroup
// 	for url := range incomingURLs() {
// 		n.Add(1)
// 		go func(url string) {
// 			defer n.Done()
// 			start := time.Now()
// 			value, err := m.Get(url, done)
// 			if err != nil {
// 				log.Print(err)
// 				return
// 			}
// 			fmt.Printf("%s, %s, %d bytes\n",
// 				url, time.Since(start), len(value.([]byte)))
// 		}(url)
// 	}
// 	n.Wait()
// 	close(done)
// }

// func httpGetBody(url string, done <-chan struct{}) (interface{}, error) {
// 	req, err := http.NewRequest("GET", url, nil)
// 	req.Cancel = done
// 	resp, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	return ioutil.ReadAll(resp.Body)
// }

// func incomingURLs() <-chan string {
// 	ch := make(chan string)
// 	go func() {
// 		for _, url := range []string{
// 			"https://golang.org",
// 			"https://godoc.org",
// 			"https://play.golang.org",
// 			"http://gopl.io",
// 			"https://golang.org",
// 			"https://godoc.org",
// 			"https://play.golang.org",
// 			"http://gopl.io",
// 		} {
// 			ch <- url
// 		}
// 		close(ch)
// 	}()
// 	return ch
// }
