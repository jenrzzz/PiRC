package pirc

type CodePair struct {
    Code int
    Msg string
}

func (p CodePair) Error() string {
    return p.Msg
}

type ProtoReply struct {
    WELCOME CodePair
    YOURHOST CodePair
    CREATED CodePair
    MYINFO CodePair
}

type ProtoError struct {
    NEEDMOREPARAMS CodePair
    NICKNAMEINUSE CodePair
}

var RPL = ProtoReply {
    WELCOME: CodePair{1, "Welcome to the Internet Relay Chat Network %v!%v@%v"},
    YOURHOST: CodePair{2, "Your host is %v, running version %v"},
    CREATED: CodePair{3, "This server was created %v"},
    MYINFO: CodePair{4, "%v %v %v %V"},
}

var ERR = ProtoError {
    NICKNAMEINUSE: CodePair{433, "Nickname is already in use"},
    NEEDMOREPARAMS: CodePair{461, "%v: Not enough parameters"},
}
