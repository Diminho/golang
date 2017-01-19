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
	conn     net.Conn
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
		// fmt.Println("Message: " + message) //debbuging statement
		if message != "" && client.id != message[:index] {
			writeMessage(client.writer, message[index+1:]+"\n")

		}
	}

}

func (room *Room) introducing(message string) {
	for _, client := range room.clients {
		writeMessage(client.writer, message)
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
func writeMessage(writer *bufio.Writer, message string) {
	writer.WriteString(message)
	writer.Flush()
}
func readMessage(reader *bufio.Reader) string {
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)
	return message
}
func newClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	writeMessage(writer, "Enter your nickname:")
	name := readMessage(reader)
	identifier := conn.RemoteAddr().String()
	client := &Client{
		name:     name,
		id:       identifier,
		incoming: make(chan string),
		writer:   writer,
		reader:   reader,
		conn:     conn}

	return client
}

func (client *Client) listen() {
	for {
		message := readMessage(client.reader)
		action, message := parseMessage(message)
		switch action {
		case "#disconnect":
			client.incoming <- client.id + "|" + " [" + client.name + "] disconnected from chat\n"
			close(client.incoming)
			client.disconnect()
			return
		case "message":
			client.incoming <- client.id + "|" + " [" + client.name + "] " + message
		case "#help":
			writeMessage(client.writer, "#disconnect - disconnects user from chat\n	#createRoom {name} - creates room with provided name\n #enter {room name} - enters to room\n")
		case "empty":

		default:
			writeMessage(client.writer, "Unknown command. Type #help to list commands\n")
		}

	}
}

func (client *Client) disconnect() {
	client.conn.Close()
}
func parseMessage(message string) (string, string) {
	var action string
	isCommand := strings.HasPrefix(message, "#")

	if isCommand {
		index := strings.Index(message, " ")
		if index == -1 {
			action = message
		} else {
			string := strings.Split(message, " ")
			action = string[0]
		}

	} else {
		action = "message"
	}
	if message == "" {
		action = "empty"
	}
	return action, message
}

func (room *Room) join(conn net.Conn) {
	client := newClient(conn)
	writeMessage(client.writer, "Welcome, "+client.name+". You have entered: "+room.name+"\n")

	go func() {
		for {
			room.incoming <- <-client.incoming
		}
	}()
	go client.listen()
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
		defer close(room.connJoins)
		defer close(room.incoming)
	}
}
