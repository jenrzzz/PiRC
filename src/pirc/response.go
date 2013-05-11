package pirc

import (
    "fmt"
)

type ServerResponse interface {
    Response(s *Server, cmd string) string
}

type CodePair struct {
    Code int
    Msg string
}

// error() interface
func (p CodePair) Error() string {
    return p.Msg
}

// ServerResponse interface
func (cp CodePair) Response(s *Server, cmd string) string {
    response := fmt.Sprintf(":%v %3d0 %v :%v\r\n", s.Hostname, cp.Code, cmd, cp.Msg)
    return response
}

// For replies that need to be formatted with info
func (cp CodePair) Format(vals ...interface{}) {
    cp.Msg = fmt.Sprintf(cp.Msg, vals)
}

// Normal server responses
var RPL = map[string] CodePair {
    "WELCOME": CodePair{1, "Welcome to the Internet Relay Chat Network %v!%v@%v"},
    "YOURHOST": CodePair{2, "Your host is %v, running version %v"},
    "CREATED": CodePair{3, "This server was created %v"},
    "MYINFO": CodePair{4, "%v %v %v %v"},
}

// error codes and messages
var ERR = map[string] CodePair {
    "UNKNOWNCOMMAND": CodePair{421, "Unknown command"},
    "NONICKNAMEGIVEN": CodePair{431, "No nickname given"},
    "ERRONEUSNICKNAME": CodePair{432, "Erroneous nickname"},
    "NICKNAMEINUSE": CodePair{433, "Nickname is already in use"},
    "NEEDMOREPARAMS": CodePair{461, "Not enough parameters"},
    "ALREADYREGISTRED": CodePair{462, "You may not reregister"},
}
