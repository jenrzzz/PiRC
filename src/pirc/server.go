package pirc

import (
    "fmt"
    "log"
    "net"
    "bytes"
)

type Server struct {
    users map[string] *User
}

func (s *Server) FindUser(nick string) *User {
    return s.users[nick]
}

func (s *Server) AddUser(u *User) error {
     if s.users[u.Nick] != nil {
        e := NickError{Code: "in_use"}
        return e
    }

    s.users[u.Nick] = u
    return nil
}

func (s *Server) AddUserByNick(nick string, conn *net.Conn) error {
    u := User{nick, conn}
    return s.AddUser(&u)
}

func RunServer(listenaddr string) {
    s := new(Server)
    ln, err := net.Listen("tcp", listenaddr)
    if err != nil {
        log.Println(err)
        log.Panicf("Unable to start a server on %v!", listenaddr)
    }

    // Main connection loop
    for {
        log.Println("Waiting for connection")
        conn, err := ln.Accept()
        if err != nil {
            log.Println("Error opening connection.")
            log.Println(err)
            continue
        } else {
            log.Println("Connection opened")
        }

        go func(c net.Conn) {
            var buf [1024]byte
            _, err := c.Read(buf[0:])
            if err != nil {
                log.Println("Error receiving data.")
                log.Println(err)
            } else {
                fmt.Printf("Received input from %v\n", c.RemoteAddr().String())
                index := bytes.IndexAny(buf[0:], "\n")
                cmd := buf[0:index]
                fmt.Println("Command: " + string(cmd))
                u, _ := CreateUser(string(cmd), &c)
                s.AddUser(u)
            }
        }(conn)
    }
}

type NickError struct {
    Code string
}

func (e NickError) Error() string {
    return "NickError: " + e.Code
}

type CmdError struct {
    Code string
    Cmd string
}

func (e CmdError) Error() string {
    return "Bad command: " + e.Cmd + "(" + e.Code + ")"
}
