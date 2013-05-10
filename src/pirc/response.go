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

// Normal server responses
type ProtoReply struct {
    WELCOME CodePair
    YOURHOST CodePair
    CREATED CodePair
    MYINFO CodePair
}

var RPL = ProtoReply {
    WELCOME: CodePair{1, "Welcome to the Internet Relay Chat Network %v!%v@%v"},
    YOURHOST: CodePair{2, "Your host is %v, running version %v"},
    CREATED: CodePair{3, "This server was created %v"},
    MYINFO: CodePair{4, "%v %v %v %V"},
}

// Error server responses
// There's probably a better way to implement these -- I'm bad at Go
type ProtoError struct {
    NEEDMOREPARAMS CodePair
    NICKNAMEINUSE CodePair
    UNKNOWNCOMMAND CodePair
}

// Define the actual error codes and messages
var ERR = ProtoError {
    UNKNOWNCOMMAND: CodePair{421, "Unknown command"},
    NICKNAMEINUSE: CodePair{433, "Nickname is already in use"},
    NEEDMOREPARAMS: CodePair{461, "Not enough parameters"},
}
