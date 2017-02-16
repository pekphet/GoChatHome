package chathome

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	cmd	string
	arg	string
}

type Run func(server *Server, client *Client, arg string)

const(
	CMD_REGEX = `:(?P<cmd>\w+)\s*(?P<arg>.*)`
)

var(
	commands map[string] Run
)

func init() {
	commands = map[string]Run{
		P_CALL_REG		:doReg,
		P_CALL_LOGIN 		:doLogin,
		P_CALL_QUIT		:doQuit,
		P_CALL_USER_LIST	:doShowUsers,
	}
}

func parseCmd(msg string) (cmd Command, err error) {
	tmp := strings.Split(msg, P_SP)
	if len(tmp) != 2 {
		err = errors.New("Command is WRONG!")
	}
	cmd.cmd = tmp[0]
	cmd.arg = tmp[1]
	return
}

func (self *Server) executeCmd(client *Client, cmd Command) (err error) {
	if f, ok := commands[cmd.cmd]; ok {
		f(self, client, cmd.arg)
		return
	}
	err = errors.New("Unknown Command: " + cmd.cmd)
	return
}

/**
 *   call methods
 */
func doReg(server *Server, client *Client, arg string) {

}

func doLogin(server *Server, client *Client, arg string) {

}

func changeName(server *Server, client *Client, arg string) {
	oldname := client.name
	client.name = arg
	server.broadcast(fmt.Sprintf("Notification: %s changed its name to %s", oldname, arg))
}

func doQuit(server *Server, client *Client, arg string) {
	client.quit()
	server.broadcast(fmt.Sprintf("Notification: %s quit the chat room.", client.name))
}

func doShowUsers(server *Server, client *Client, arg string) {
	var users string = "Command Result:user_list<-ALL"
	for _, r_client := range server.clients {
		users = users + "," + r_client.name
	}
	client.outgoing <- users
}

