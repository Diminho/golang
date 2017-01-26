package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

const (
	CONN_TYPE = "tcp"
	LOBBY     = "main"
)

type Client struct {
	id       string
	name     string
	incoming chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
	conn     net.Conn
	room     string
}

type Room struct {
	clients   map[string]*Client
	name      string
	connJoins chan net.Conn
	incoming  chan string
	messages  int64
}

var allClients map[string]*Client
var allRooms map[string]*Room
var messages int64
var ticker *time.Ticker

func getClientByConn(id string) *Client {
	client := allClients[id]
	return client
}
func getRoomByName(name string) (*Room, string) {
	var error string
	room, present := allRooms[name]
	if !present {
		error = "No such room"
	}
	return room, error
}

func (room *Room) post(message string) {

	for _, client := range room.clients {
		index := strings.Index(message, "|")
		if message != "" && client.id != message[:index] {
			writeMessage(client.writer, message[index+1:]+"\r\n")

		} else if index == -1 {
			room.printMetrics(message)
		}
	}

}

func (room *Room) printMetrics(message string) {
	for _, client := range room.clients {
		writeMessage(client.writer, "[Announcer] "+message+"\r\n")
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
				atomic.AddInt64(&room.messages, 1)
			case conn := <-room.connJoins:
				client, present := allClients[conn.RemoteAddr().String()]
				if !present {
					client = newClient(conn)
				}
				room.join(client)
				// case <-ticker.C:
				// 	roomMessCount := fmt.Sprintf("%d", room.messages)
				// 	room.incoming <- "Messages sent in room: " + roomMessCount
				// 	// room.printMetrics("Messages sent in room: " + roomMessCount)
				// 	fmt.Println("HELLO TICKER")
			}
		}
	}()
}
func writeMessage(writer *bufio.Writer, message string) {
	writer.WriteString(message)
	writer.Flush()
}
func readMessage(reader *bufio.Reader) (string, error) {
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	message = strings.TrimSpace(message)
	return message, nil
}
func newClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	writeMessage(writer, "Enter your nickname: ")
	name, _ := readMessage(reader)
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

func (client *Client) listen(room *Room) {
	var action string
	var message string
	var error error
	for {
		message, error = readMessage(client.reader)
		if error != nil {
			if error.Error() == "EOF" {
				action = "#disconnect"
			}
		} else {
			action, message = parseMessage(message)
		}

		switch action {
		case "#disconnect":
			client.disconnect(room)
			return
		case "#createRoom":
			customRoom := createRoom(message)
			customRoom.connJoins <- client.conn
			client.leaveRoom(room)
			delete(room.clients, client.id)
			return
		case "#enterRoom":
			customRoom, error := getRoomByName(message)
			if error == "" {
				fmt.Println(error)
			}
			customRoom.connJoins <- client.conn
			client.leaveRoom(room)
			delete(room.clients, client.id)
		case "#leave":
			client.leaveRoom(room)
			delete(room.clients, client.id)
			lobbyRoom, error := getRoomByName(LOBBY)
			if error == "" {
				fmt.Println(error)
			}
			lobbyRoom.connJoins <- client.conn

		case "message":
			client.incoming <- client.id + "|" + " [" + client.name + "] " + message
		case "#help":
			writeMessage(client.writer, " #disconnect - disconnects user from chat\r\n #createRoom {name} - creates room with provided name\r\n #enterRoom {room name} - enters to room\r\n #leave - user leaves current room and goes to main lobby")
		case "empty":

		default:
			writeMessage(client.writer, "Unknown command. Type #help to list commands\r\n")
		}

	}
}

func (client *Client) leaveRoom(room *Room) {
	delete(room.clients, client.id)
	client.incoming <- client.id + "|" + " [" + client.name + "] left the room " + room.name + "\r\n"
}

func (client *Client) disconnect(room *Room) {
	client.incoming <- client.id + "|" + " [" + client.name + "] disconnected from chat\r\n"
	delete(room.clients, client.id)
	delete(allClients, client.id)
	close(client.incoming)
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
			message = string[1]
		}

	} else {
		action = "message"
	}
	if message == "" {
		action = "empty"
	}
	return action, message
}

func (room *Room) join(client *Client) {
	writeMessage(client.writer, "Welcome, "+client.name+". You have entered: "+room.name+"\r\n")

	go func() {
		for {
			message, isClosed := <-client.incoming
			if isClosed {
				room.incoming <- message
			}
			_, present := room.clients[client.id]
			if !present {
				return
			}
		}
	}()
	go client.listen(room)
	room.introducing("User [" + client.name + "] entered the chat\r\n")
	client.room = room.name
	room.clients[client.id] = client
	allClients[client.id] = client
}

func createRoom(name string) *Room {
	joinsChan := make(chan net.Conn)
	room := &Room{
		name:      name,
		connJoins: joinsChan,
		clients:   map[string]*Client{},
		incoming:  make(chan string)}
	allRooms[room.name] = room
	room.startRoomChat()

	return room
}
func init() {
	allClients = map[string]*Client{}
	allRooms = map[string]*Room{}
	ticker = time.NewTicker(time.Second * 30)
}

func main() {

	room := createRoom(LOBBY)
	listener, _ := net.Listen(CONN_TYPE, "localhost:8989")
	defer listener.Close()
	for {
		connection, _ := listener.Accept()
		room.connJoins <- connection
	}

}
