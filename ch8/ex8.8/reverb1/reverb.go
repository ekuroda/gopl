package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	inputChan := make(chan string)
	doneChan := make(chan struct{})

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)

	go func() {
		for input.Scan() {
			if timer.Stop() {
				timer.Reset(timeout)
				inputChan <- input.Text()
			}
		}
		close(inputChan)
		fmt.Printf("input scan done\n")
		doneChan <- struct{}{}
	}()

	go func() {
		timeout := false
		for !timeout {
			select {
			case text, ok := <-inputChan:
				if ok {
					echo(c, text, 1*time.Second)
				}
			case <-timer.C:
				timeout = true
			default:
			}
		}
		c.Close()
	}()

	<-doneChan
	fmt.Printf("done\n")
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
