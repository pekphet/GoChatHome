package chathome

import (
	"log"
	"strings"
	"net"
	_ "chathome/db"
	"strconv"
)

const (
	MAX_CLIENTS = 100
)

func (self *Server) generateToken() {
	self.tokens <- 0
}

func (self *Server) takeToken() {
	<-self.tokens
}

func CreateServer() *Server {
	server := &Server{
		connClients:	make(ConnClientTable, MAX_CLIENTS),
		clients:        make(ClientTable, MAX_CLIENTS),
		tokens:         make(Token, MAX_CLIENTS),
		pending:        make(chan net.Conn),
		quiting:        make(chan net.Conn),
		incoming:       make(Message),
		outgoing:       make(Message),
		broadcasting:   make(Message),
	}
	server.listen()
	return server
}

func (self *Server) listen() {
	go func() {
		for {
			select {
			case msg := <-self.incoming:
				self.sending(msg)
			case msg := <-self.broadcasting:
				self.broadcast(msg)
			case conn := <-self.pending:
				self.join(conn)
			case conn := <-self.quiting:
				self.leave(conn)
			}

		}
	}()
}

func (self *Server) join(conn net.Conn) {
	client := CreateClient(conn)
	name := getUniqName()
	client.name = name
	self.connClients[conn] = client
	log.Printf("Auto assigned name for conn %p: %s\n", conn, name)
	go func() {
		for {
			msg := <-client.incoming
			log.Printf("Got[%s] from client %s\n", msg, client.name)
			if (!strings.Contains(msg, P_SP)) {
				continue
			}
			if strings.HasPrefix(msg, P_SEND_MSG) {
				if strings.HasPrefix(msg, P_SEND_MSG + P_SP) {
					self.broadcasting <- makeMsg(client.name, msg)
				} else if strings.HasPrefix(msg, P_SEND_MSG + P_SP_SEND) {
					self.incoming <- makeSendingMsg(client.name, msg)
				}
			} else if strings.HasPrefix(msg, P_SEND_G_MSG + P_SP) {
				//TODO DEVELOP LATER
			} else {
				if cmd, err := parseCmd(msg); err == nil {
					if err = self.executeCmd(client, cmd); err == nil {
						continue
					} else {
						log.Println(err.Error())
					}
				} else {
					log.Println(err.Error())
				}
			}
		}
	}()

	go func() {
		for {
			conn := <-client.quiting
			log.Printf("Client %s is quiting\n", client.name)
			self.quiting <- conn
		}
	}()
}

func (self *Server) broadcast(message string) {
	log.Printf("broadcast: %s\n", message)
	for _, client := range self.clients {
		client.incoming <- message
	}
}

func (self *Server) sending(message string) {
	cms := strings.Split(message, P_SP_RCV)
	uid, err := strconv.Atoi(cms[0])
	if err != nil {
		return
	}
	self.clients[uid].incoming <- cms[1]
}


func (self *Server) leave(conn	net.Conn) {
	if conn != nil {
		delete(self.clients, self.connClients[conn].uid)
		conn.Close()
		delete(self.connClients, conn)
	}
	self.generateToken()
}

func (self *Server) Start(connStr string) {
	self.listener, _ = net.Listen("tcp", connStr)
	log.Printf("Server %p starts\n", self)

	for i := 0; i < MAX_CLIENTS; i++ {
		self.generateToken()
	}

	for {
		conn, err := self.listener.Accept()

		if err != nil {
			log.Println(err)
			self.leave(conn)
			return
		}
		log.Printf("A new conn %v kicks\n", conn)

		self.takeToken()
		self.pending <- conn
	}
}
func (self *Server) Stop() {
	self.listener.Close()
}

func (server *Server) makeClientUIDIndex(c *Client) {
	server.clients[c.uid] = c
}

/***SMART METHODS***/
func makeMsg(name string, msg string) string {
	return P_SEND_MSG + P_SP + name + ":" + strings.Split(msg, P_SP)[1]
}

func makeSendingMsg(name string, msg string) string {
	cms := strings.Split(msg, P_SP)
	mms := strings.Split(cms[0], P_SP_SEND)
	return mms[1] + P_SP_RCV + P_SEND_MSG + P_SP_SEND +  P_SP + name + ":" + cms[1]
}
