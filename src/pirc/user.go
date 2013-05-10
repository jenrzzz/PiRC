package pirc

import (
    "net"
    "strings"
)

type User struct {
    Nick string
    Conn *net.Conn
}

func CreateUser(cmd string, c *net.Conn) (*User, error) {
    index := strings.LastIndex(cmd, " ")
    // if strings.ToUpper(cmd[:index]) != "NICK" {
    //     return nil, CmdError{"bad", cmd[:index]}
    // }

    u := User{cmd[index+1:], c}
    return &u, nil
}
