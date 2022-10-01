package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/wagaru/go-server/http"

	log "github.com/sirupsen/logrus"
)

var SERVER_PORT = ":8080"

const (
	DEFAULT_TIMEOUT = 30 * time.Second
)

func main() {
	server, err := net.Listen("tcp", SERVER_PORT)
	if err != nil {
		log.Panicf("cannot start server on %v", SERVER_PORT)
	}
	defer server.Close()

	log.Printf("start listening on %v ...\n", SERVER_PORT)
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("conn has error: %v\n", err)
			continue
		}
		go handleFunc(conn)
	}
}

func handleFunc(conn net.Conn) {
	defer conn.Close()

	// read & write timeout after DEFAULT_TIMEOUT
	conn.SetDeadline(time.Now().Add(DEFAULT_TIMEOUT))

	r := bufio.NewReader(conn)
	requestChan := make(chan string, 10)
	for {
		bytes, err := r.ReadBytes('\n')
		if err != nil {
			// error will be EOF when conn closed
			log.Infof("err: %v", err)
			close(requestChan)
			break
		}
		if length := len(bytes); length > 1 && bytes[length-2] == '\r' {
			requestChan <- string(bytes[:len(bytes)-2])
		}
	}

	res := make(chan http.Request)
	go func() {
		httpRequest := &http.Request{}
		for data := range requestChan {
			err := http.ParseReceivedData(data, httpRequest)
			if err != nil {
				panic("TBD")
			}
		}
		res <- *httpRequest
	}()

	d := <-res
	fmt.Printf("%+v", d)
	for _, v := range d.Headers {
		fmt.Printf("%v: %v", v.Name, v.Value)
	}

	responseBody := "Hello World!"
	response := []string{
		"HTTP/1.1 200 OK",
		"Date: " + time.Now().Format(time.RFC1123),
		"",
		responseBody,
	}

	_, err := conn.Write([]byte(strings.Join(response, "\r\n")))
	if err != nil {
		log.Printf("write content has error: %v", err)
	}
}
