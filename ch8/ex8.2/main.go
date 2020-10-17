package main

// https://tools.ietf.org/html/rfc959
// http://srgia.com/docs/rfc959j.html

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const commandPort = 8020
const dataPort = 8021

type clientHandler struct {
	conn         net.Conn
	dataListener net.Listener
	cd           string
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", commandPort))
	if err != nil {
		log.Fatal(err)
	}

	dataListener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", dataPort))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleCommand(conn, dataListener)
	}
}

func handleCommand(c net.Conn, dataListener net.Listener) {
	log.Printf("connected")
	writeReply(c, 200, "")

	ch := clientHandler{conn: c, dataListener: dataListener, cd: ""}

	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		command, params, err := parseCommand(scanner.Text())
		if err != nil {
			writeReply(c, 550, "")
			continue
		}

		switch command {
		case "USER":
			// 誰でもok
			writeReply(c, 200, "")
		case "SYST":
			// USERに続いて来るので返しておく
			writeReply(c, 200, "")
		case "FEAT":
			// USERに続いて来るので返しておく
			writeReply(c, 200, "")
		case "EPSV":
			// EPSV非対応。PASVが来る想定
			writeReply(c, 550, "")
		case "PASV":
			ch.handlePASV(params)
		case "TYPE":
			ch.handleTYPE(params)
		case "CWD":
			ch.handleCWD(params)
		case "LIST":
			ch.handleLIST()
		case "RETR":
			ch.handleRETR(params)
		case "QUIT":
			ch.handleQUIT()
			break
		default:
			// 対応してないコマンド
			writeReply(c, 502, "")
		}
	}
	log.Printf("disconnected")
}

func writeReply(c net.Conn, code int, message string) {
	if message == "" {
		log.Printf("%d", code)
		io.WriteString(c, fmt.Sprintf("%d\r\n", code))
	} else {
		log.Printf("%d %s", code, message)
		io.WriteString(c, fmt.Sprintf("%d %s\r\n", code, message))
	}
}

func parseCommand(line string) (command, params string, err error) {
	log.Printf("command: %s", line)
	splitted := strings.SplitN(line, " ", 2)
	if len(splitted) == 0 {
		err = fmt.Errorf("invalid command")
		return
	}

	command = splitted[0]
	if len(splitted) == 2 {
		params = splitted[1]
	}
	return
}

func getIPAddressQuads(c net.Conn) []string {
	ip := strings.Split(c.LocalAddr().String(), ":")[0]
	quads := strings.Split(ip, ".")
	return quads
}

func (ch *clientHandler) handlePASV(param string) {
	quads := getIPAddressQuads(ch.conn)
	port := []int{dataPort >> 8, dataPort & 255}
	//log.Printf("quads=%v, port=%v", quads, port)
	writeReply(ch.conn, 227, fmt.Sprintf("Entering Passive Mode (%s,%d,%d)", strings.Join(quads, ","), port[0], port[1]))
}

func (ch *clientHandler) handleTYPE(param string) {
	// TYPE
	// 200
	// 500, 501, 504, 421, 530
	if param != "I" {
		log.Printf("invalid type: %s", param)
		writeReply(ch.conn, 503, "")
		return
	}
	writeReply(ch.conn, 200, "")
}

func (ch *clientHandler) handleCWD(param string) {
	// CWD
	// 250
	// 500, 501, 502, 421, 530, 550
	cd := ch.cd
	if param == "" {
		log.Print("directory required")
		writeReply(ch.conn, 550, "")
		return
	}

	cd = filepath.Join(cd, param)
	path := filepath.Join(".", cd)

	file, err := os.Stat(path)
	if err != nil {
		log.Printf("file not exists: path=%s, err=%s", path, err)
		writeReply(ch.conn, 550, "")
		return
	}
	if !file.IsDir() {
		log.Printf("not dir: path=%s", path)
		writeReply(ch.conn, 550, "")
		return
	}

	ch.cd = cd
	writeReply(ch.conn, 250, fmt.Sprintf("directory changed to %s", cd))
}

func (ch *clientHandler) handleLIST() {
	// LIST
	// 125, 150
	//    226, 250
	//    425, 426, 451
	// 450
	// 500, 501, 502, 421, 530
	path := "."
	if ch.cd != "" {
		path = filepath.Join(path, ch.cd)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("failed to read dir: path=%s, err=%s", path, err)
		writeReply(ch.conn, 550, "")
		return
	}

	log.Printf("path=%s", path)

	writeReply(ch.conn, 150, "")
	conn, err := ch.dataListener.Accept()
	if err != nil {
		log.Print(err)
		writeReply(ch.conn, 425, "")
		return
	}

	defer conn.Close()

	var paths []string
	for _, file := range files {
		log.Printf("> %s", file.Name())
		paths = append(paths, file.Name())
	}

	_, err = fmt.Fprintf(conn, strings.Join(paths, "\r\n"))
	if err != nil {
		log.Print(err)
		writeReply(ch.conn, 451, "")
	}

	writeReply(ch.conn, 226, "")
}

func (ch *clientHandler) handleRETR(params string) {
	// RETR
	// 125, 150
	//    (110)
	//    226, 250
	//    425, 426, 451
	// 450, 550
	// 500, 501, 421, 530
	path := filepath.Join(".", ch.cd, params)

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Printf("file not exists: path=%s, err=%s", path, err)
		writeReply(ch.conn, 550, "")
		return
	}
	if fileInfo.IsDir() {
		log.Printf("not file: path=%s", path)
		writeReply(ch.conn, 550, "")
		return
	}

	log.Printf("path=%s", path)

	writeReply(ch.conn, 150, "")
	conn, err := ch.dataListener.Accept()
	if err != nil {
		log.Print(err)
		writeReply(ch.conn, 425, "")
		return
	}

	defer conn.Close()

	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		writeReply(ch.conn, 451, "")
		return
	}

	defer file.Close()

	io.Copy(conn, file)
	if err != nil {
		if err != io.EOF {
			log.Print(err)
			writeReply(ch.conn, 451, "")
			return
		}
	}

	writeReply(ch.conn, 226, "")
}

func (ch *clientHandler) handleQUIT() {
	writeReply(ch.conn, 221, "")
	ch.conn.Close()
}
