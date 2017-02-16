package chathome

import (
	"log"
	"strings"
	"net"
	"chathome/db"
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
		clients:        make(ClientTable, MAX_CLIENTS),
		tokens:         make(Token, MAX_CLIENTS),
		pending:        make(chan net.Conn),
		quiting:        make(chan net.Conn),
		incoming:       make(Message),
		outgoing:       make(Message),
		broadcasting:   make(Message),
	}
	db.StartDB()
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
		client.outgoing <- message
	}
}

func (self *Server) sending(message string) {

}


func (self *Server) leave(uid int) {
	if self.clients[uid].conn != nil {
		delete(self.connClients, self.clients[uid].conn)
		self.clients[uid].conn.Close()
		delete(self.clients, uid)
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

/***SMART METHODS***/
func makeMsg(name string, msg string) string {
	return P_SEND_MSG + P_SP + name + ":" + strings.Split(msg, P_SP)[1]
}

func makeSendingMsg(name string, msg string) string {
	P_SEND_MSG + P_SP_RCV
}


func getUserFromMsg(msg string) string {
	return strings.Split(msg, ":")[0]
}

func parseSentMsg(msg string) (string, string) {
	tmpStr := strings.Split(msg, ":")
	log.Printf("tmpstr0:%s, 1:%s", tmpStr[0], tmpStr[1])
	tmpp := strings.Split(tmpStr[0], "->")
	log.Printf("tmpp0:%s, 1:%s", tmpp[0], tmpp[1])
	return tmpp[0] + ":" + tmpStr[1], tmpp[1]
}