package main

import (
	"bufio"
	"net"
	"strings"
)

const (
	CONN_TYPE = "tcp"
)

type Client struct {
	id       string
	name     string
	incoming chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

type Room struct {
	clients   map[string]*Client
	name      string
	connJoins chan net.Conn
	incoming  chan string
}

func (room *Room) post(message string) {
	for _, client := range room.clients {
		index := strings.Index(message, "|")
		if client.id != message[:index] {
			client.writer.WriteString(message[index+1:])
			client.writer.Flush()
		}
	}

}

func (room *Room) introducing(message string) {
	for _, client := range room.clients {
		client.writer.WriteString(message)
		client.writer.Flush()
	}
}

func (room *Room) startRoomChat() {
	go func() {
		for {
			select {
			case message := <-room.incoming:
				room.post(message)
			case conn := <-room.connJoins:
				room.join(conn)
			}
		}
	}()
}
func newClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	identifier := conn.RemoteAddr().String()
	client := &Client{
		name:     name[:len(name)-2], //triming \n chars
		id:       identifier,
		incoming: make(chan string),
		writer:   writer,
		reader:   reader}

	return client
}

func (client *Client) listen() {
	go func() {
		for {
			message, _ := client.reader.ReadString('\n')
			client.incoming <- client.id + "|" + " [" + client.name + "] " + message
		}
	}()
}

func (room *Room) join(conn net.Conn) {
	client := newClient(conn)
	client.writer.WriteString("Welcome, " + client.name + ". You have entered: " + room.name + "\n")
	client.writer.Flush()

	go func() {
		for {
			room.incoming <- <-client.incoming

		}

	}()
	client.listen()
	room.introducing("User [" + client.name + "] entered the chat\n")
	room.clients[client.id] = client
}

func createRoom(name string) *Room {
	joinsChan := make(chan net.Conn)
	room := &Room{
		name:      name,
		connJoins: joinsChan,
		clients:   map[string]*Client{},
		incoming:  make(chan string)}

	room.startRoomChat()

	return room
}

func main() {

	room := createRoom("lobby")
	listener, _ := net.Listen(CONN_TYPE, "localhost:8989")
	defer listener.Close()
	for {
		connection, _ := listener.Accept()
		room.connJoins <- connection

	}
}
