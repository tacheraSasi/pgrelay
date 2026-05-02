package main

import (
	"io"
	"log"
	"net"
)

const (
	LISTEN_ADDR = ":5433"
	PG_ADDR     = "localhost:5432"
)

func main() {
	ln, err := net.Listen("tcp", LISTEN_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("proxy listening on", LISTEN_ADDR)

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handle(clientConn)
	}
}

func handle(client net.Conn) {
	defer client.Close()

	server, err := net.Dial("tcp", PG_ADDR)
	if err != nil {
		log.Println(err)
		return
	}

	defer server.Close()

	go io.Copy(server, client)
	io.Copy(client, server)
}
