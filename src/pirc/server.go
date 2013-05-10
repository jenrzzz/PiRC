package pirc

import (
    "fmt"
    "log"
    "net"
    "bytes"
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

func RunServer(listenaddr string) {
    s := new(Server)
    s.users = make(map[string] *User)
    s.Hostname = "localhost"
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
            bytes_read, err := c.Read(buf[0:])
            if err != nil {
                log.Println("Error receiving data.")
                log.Println(err)
            } else {
                fmt.Printf("Received input of length %d from %v\n", bytes_read, c.RemoteAddr().String())
                cmd_term := bytes.IndexAny(buf[0:], "\n")
                last_term := 0
                for cmd_term != -1 && last_term < bytes_read {
                    fmt.Printf("Accessing buf[%d:%d]", last_term, cmd_term)
                    cmd := buf[last_term:cmd_term]
                    fmt.Println("Command: " + string(cmd))
                    fmt.Println("First: " + string(cmd)[:bytes.IndexAny(cmd, " ")])

                    if string(cmd)[:bytes.IndexAny(cmd, " ")] == "NICK" {
                        var err *CodePair = nil
                        var u *User
                        u, _ = CreateUser(string(cmd), &c)
                        err = s.AddUser(u)

                        if err != nil {
                            WriteResponse("nick", err, c, s)
                        } else {
                            log.Println("Created user " + u.Nick)
                        }
                    }

                    last_term = cmd_term + 1
                    cmd_term = bytes.IndexAny(buf[last_term:], "\n") + last_term
                }
            }
        }(conn)
    }
}

func WriteResponse(cmd string, cp *CodePair, c net.Conn, s *Server) {
    response := fmt.Sprintf(":%v %3d0 %v :%v\r\n", s.Hostname, cp.Code, cmd, cp.Msg)
    c.Write([]byte(response))
}

