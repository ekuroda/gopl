package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
)

type clock struct {
	loc  string
	addr string
	conn net.Conn
	time string
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("no arguments")
	}

	clocks := make([]*clock, 0)
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		locAndAddr := strings.Split(arg, "=")
		if len(locAndAddr) != 2 {
			log.Printf("ignore invalid arg :%q\n", arg)
			continue
		}
		conn, err := net.Dial("tcp", locAndAddr[1])
		if err != nil {
			log.Printf("skip location: %q; failed to dial: %q\n", locAndAddr[0], locAndAddr[1])
			continue
		}
		clocks = append(clocks, &clock{locAndAddr[0], locAndAddr[1], conn, ""})
	}

	defer func() {
		for _, clock := range clocks {
			clock.conn.Close()
		}
	}()

	var wg sync.WaitGroup
	update := make(chan struct{})
	for _, clock := range clocks {
		wg.Add(1)
		go updateClock(clock, update, &wg)
	}

	go func() {
		wg.Wait()
		close(update)
	}()

	for range update {
		printClocks(clocks)
	}
}

func updateClock(clock *clock, update chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(clock.conn)
	for scanner.Scan() {
		clock.time = scanner.Text()
		update <- struct{}{}
	}
	clock.time = "--:--:--"
	update <- struct{}{}
}

func printClocks(clocks []*clock) {
	const format = "%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 10, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Location", "Time")
	fmt.Fprintf(tw, format, "--------", "--------")
	for _, clock := range clocks {
		fmt.Fprintf(tw, format, clock.loc, clock.time)
	}
	// https://stackoverflow.com/questions/15442292/how-to-have-an-in-place-string-that-updates-on-stdout
	// https://qiita.com/PruneMazui/items/8a023347772620025ad6
	fmt.Fprintf(tw, "\033[%dA", len(clocks)+2)
	tw.Flush()
}
