package chathome

import (
	"net"
	"bufio"
	"log"
)

type Message chan string

type Client struct {
	uid      int
	name     string
	token    string
	conn     net.Conn
	incoming Message
	outgoing Message
	reader   *bufio.Reader
	writer   *bufio.Writer
	quiting  chan net.Conn
}

func CreateClient(conn net.Conn) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	client := &Client{
		uid	: 0,
		conn	: conn,
		incoming: make(Message),
		outgoing: make(Message),
		quiting	: make(chan net.Conn),
		reader	: reader,
		writer	: writer,
	}
	client.Listen()
	return client
}

func (client *Client) Listen() {
	go client.read()
	go client.write()
}

func (client *Client) read() {
	for {
		if line, _, err := client.reader.ReadLine(); err == nil {
			client.incoming <- string(line)
		} else {
			log.Printf("Read error: %s\n", err)
			client.quit()
			return
		}
	}
}

func (client *Client) write() {
	for data := range client.outgoing {
		if _, err := client.writer.WriteString(data + "\n"); err != nil {
			client.quit()
			return
		}

		if err := client.writer.Flush(); err != nil {
			log.Printf("Write error: %s\n", err)
			client.quit()
			return
		}
	}
}

func (client *Client) quit() {
	client.quiting <- client.conn
}

func (client *Client) GetIncoming() string {
	return <- client.incoming
}

func (client *Client) PutOutgoing(msg string) {
	client.outgoing <- msg
}
