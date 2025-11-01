package main

import (
	"log"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleChatConn(conn)
	}
}

func handleChatConn(conn net.Conn) {

}

func broadcast() {
	clients := make(map[client]bool)
	select {
	case msg := <-messages:
		for c := range clients {
			c <- msg
		}
	case cli := <-entering:
		clients[cli] = true
	case cli := <-leaving:
		delete(clients, cli)
		close(cli)
	}
}
