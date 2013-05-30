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
    response := fmt.Sprintf(":%v %03d %v :%v\r\n", s.Hostname, cp.Code, cmd, cp.Msg)
    return response
}

// For replies that need to be formatted with info
func (cp CodePair) Format(vals ...interface{}) CodePair {
    ncp := new(CodePair)
    ncp.Code = cp.Code
    ncp.Msg = fmt.Sprintf(cp.Msg, vals...)
    return *ncp
}

func (cp CodePair) FormatMsg(vals ...interface{}) string {
    return fmt.Sprintf(cp.Msg, vals...)
}

// Normal server responses
var RPL = map[string] CodePair {
    "WELCOME": CodePair{1, "Welcome to the Internet Relay Chat Network %v"},
    "YOURHOST": CodePair{2, "Your host is %v, running version %v"},
    "CREATED": CodePair{3, "This server was created %v"},
    "MYINFO": CodePair{4, "%v %v %v %v"},
    "NOTOPIC": CodePair{331, "%v :No topic is set"},
    "TOPIC": CodePair{332, "%v :%v"},
    "NAMREPLY": CodePair{353, "%v :%v"},
    "ENDOFNAMES": CodePair{366, "%v :End of /NAMES list"},
    "MOTD": CodePair{375, "%v"},
    "MOTDSTART": CodePair{375, "- %v Message of the day -"},
    "ENDOFMOTD": CodePair{376, "End of /MOTD command"},
}

// error codes and messages
var ERR = map[string] CodePair {
    "NOSUCHCHANNEL": CodePair{403, "No such channel"},
    "UNKNOWNCOMMAND": CodePair{421, "Unknown command %v"},
    "NONICKNAMEGIVEN": CodePair{431, "No nickname given"},
    "ERRONEUSNICKNAME": CodePair{432, "Erroneous nickname"},
    "NICKNAMEINUSE": CodePair{433, "Nickname is already in use"},
    "NOTONCHANNEL": CodePair{442, "You're not on that channel"},
    "NOTREGISTERED": CodePair{451, "You have not registered"},
    "NEEDMOREPARAMS": CodePair{461, "Not enough parameters"},
    "ALREADYREGISTRED": CodePair{462, "You may not reregister"},
    "USERSDONTMATCH": CodePair{502, "Can't change mode for other users"},
}
