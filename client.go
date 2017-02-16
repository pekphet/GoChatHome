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
		conn:  	conn,
		incoming: make(Message),
		outgoing: make(Message),
		quiting:  make(chan net.Conn),
		reader:   reader,
		writer:   writer,
	}
	client.Listen()
	return client
}

func (self *Client) Listen() {
	go self.Read()
	go self.Write()
}

func (self *Client) Read() {
	for {
		if line, _, err := self.reader.ReadLine(); err == nil {
			self.incoming <- string(line)
		} else {
			log.Printf("Read error: %s\n", err)
			self.quit()
			return
		}
	}
}

func (self *Client) Write() {
	for data := range self.outgoing {
		if _, err := self.writer.WriteString(data + "\n"); err != nil {
			self.quit()
			return
		}

		if err := self.writer.Flush(); err != nil {
			log.Printf("Write error: %s\n", err)
			self.quit()
			return
		}
	}
}

func (self *Client) quit() {
	self.quiting <- self.conn
}

func (self *Client) GetIncoming() Message {
	return self.incoming
}