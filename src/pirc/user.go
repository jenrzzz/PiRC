package pirc

import (
    "net"
)

type User struct {
    Nick string
    Conn net.Conn
}

func CreateUser(nick string, c net.Conn) (*User, error) {
    u := User{nick, c}
    return &u, nil
}
