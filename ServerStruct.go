package chathome

import "net"

type ClientTable 	map [int] *Client
type ConnClientTable 	map [net.Conn] *Client
type Token		chan int

type Server struct {
	listener	net.Listener
	clients		ClientTable
	connClients	ConnClientTable
	tokens		Token
	pending		chan net.Conn
	quiting		chan net.Conn
	broadcasting	Message
	incoming	Message
	outgoing	Message
}

