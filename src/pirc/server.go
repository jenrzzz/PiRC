package pirc

import (
    // "fmt"
    "log"
    "net"
    // "bytes"
)

type Server struct {
    Hostname string
    users map[string] *User
    // channels map[string] *Channel
}

func (s *Server) FindUser(nick string) *User {
    return s.users[nick]
}

func (s *Server) AddUser(u *User) *CodePair {
     if s.users[u.Nick] != nil {
        return &ERR.NICKNAMEINUSE
    }

    s.users[u.Nick] = u
    return (*CodePair)(nil)
}

func (s *Server) AddUserByNick(nick string, conn *net.Conn) *CodePair {
    u := User{nick, conn}
    return s.AddUser(&u)
}

var server = new(Server)

func RunServer(listenaddr string) {
    server.users = make(map[string] *User)
    server.Hostname = "localhost"
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

        go func(conn net.Conn) {
            var buf [1024]byte
            bytes_read, err := conn.Read(buf[0:])
            if err != nil {
                log.Println("Error receiving data.")
                log.Println(err)
            } else {
                log.Printf("Received input of length %d from %v\n", bytes_read, conn.RemoteAddr().String())
                parser := new(CmdParser)
                cmd, err := parser.Parse(buf[0:bytes_read])
                if err != nil {
                    conn.Write([]byte(err.Response(server, cmd)))
                }
            }
        }(conn)
    }
}

