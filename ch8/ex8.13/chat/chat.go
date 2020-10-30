package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	name string
	c    chan<- string
}

var (
	entering = make(chan *client)
	leaving  = make(chan *client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[*client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.c <- msg
			}

		case cli := <-entering:
			cli.c <- "Members:"
			for c := range clients {
				cli.c <- "  " + c.name
			}
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.c)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	cli := &client{who, ch}
	messages <- who + " has arrived"
	entering <- cli

	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)

	leavingMessage := who + " has left"
	go func() {
		<-timer.C
		leavingMessage = who + " has kicked out"
		conn.Close()
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if timer.Stop() {
			timer.Reset(timeout)
			messages <- who + ": " + input.Text()
		}
	}

	leaving <- cli
	messages <- leavingMessage
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
