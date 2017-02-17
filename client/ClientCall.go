package client

import (
	"fmt"
	"net"
	"log"
	. "chathome"
)

var conn 	net.Conn
var client      *Client
var err 	error

func Conn() bool {
	conn, err = net.Dial("tcp", "www.lv90.cn:59000")
	if err != nil {

		log.Fatal(err)
	}
	client = CreateClient(conn)
	return err != nil
}

func Receive() string {
	return client.GetIncoming()
}

func SendMsg(msg string) {
	client.PutOutgoing(P_SEND_MSG + P_SP + msg)
}

func SendMsgTo(uid int, msg string) {
	client.PutOutgoing(fmt.Sprintf("%s%s%d%s%s", P_SEND_MSG, P_SP_SEND, uid, P_SP, msg))
}

func Quit() {
	client.PutOutgoing(P_CALL_QUIT + P_SP)
}

func List() {
	client.PutOutgoing(P_CALL_USER_LIST + P_SP)
}

func Call(method string, args string) {
	client.PutOutgoing(method + P_SP + args)
}