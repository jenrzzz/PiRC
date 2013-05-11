package pirc

import (
    "log"
    "net"
    "strings"
)

type IrcConn struct {
    rwc net.Conn
    remoteAddr string
    server *Server
    body []byte
}

// Implement Read/Write interface on conn (delegates to actual connection)
func (c *IrcConn) Read(b []byte) (n int, err error) {
    return c.rwc.Read(b)
}

func (c *IrcConn) Write(b []byte) (n int, err error) {
    return c.rwc.Write(b)
}

// Das server
type Server struct {
    Hostname string
    UsersByNick map[string] *User
    UsersByUsername map[string] *User
    Connections map[*IrcConn] *User
    // channels map[string] *Channel
}

func (s *Server) RegisterConnection(u *User) {
    s.Connections[u.Conn] = u
}

func (s *Server) RegisterUser(u *User) {
    s.UsersByNick[u.Nick] = u
    s.UsersByUsername[u.Username] = u
    s.Connections[u.Conn] = u
    u.Registered = true
}

func (s *Server) ChangeNick(u *User, nick string) ServerResponse {
    if s.UsersByNick[nick] != nil {
        return ERR["NICKNAMEINUSE"]
    }

    log.Printf("User %v changed nick from %v to %v", u.Username, u.Nick, nick)

    // Delete the old nick if we're changing it
    delete(server.UsersByNick, u.Nick)
    u.Nick = nick
    s.UsersByNick[u.Nick] = u
    return nil
}

// Finders
func (s *Server) FindUserByName(name string) *User {
    if u, exists := s.UsersByUsername[name]; exists {
        return u
    }
    return nil
}

func (s *Server) FindUserByNick(nick string) *User {
    if u, exists := s.UsersByNick[nick]; exists {
        return u
    }
    return nil
}

func (s *Server) FindUserByConn(c *IrcConn) *User {
    if u, exists := s.Connections[c]; exists {
        return u
    }
    return nil
}

// Check if command is valid and dispatch, then write response
func (s *Server) Dispatch(cmds *CmdParser) {
    for c := cmds.Next(); c != nil; c = cmds.Next() {
        log.Printf("Dispatching command %v", *c)
        var response ServerResponse
        if cmd_func, exists := CmdDispatcher[c.Cmd]; exists {
            response = cmd_func(c, cmds.Client)
        } else {
            response = ERR["UNKNOWNCOMMAND"]
        }

        if response != nil {
            cmds.Client.Write([]byte(response.Response(s, c.Cmd)))
        }
    }
}

var server = new(Server)
func RunServer(listenaddr string) {
    server.UsersByNick = make(map[string] *User)
    server.UsersByUsername = make(map[string] *User)
    server.Connections = make(map[*IrcConn] *User)
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
            remote_addr := conn.RemoteAddr().String()
            remote_addr = remote_addr[0:strings.IndexAny(remote_addr, ":")]
            c := IrcConn{conn, remote_addr, server, make([]byte, 2048)}
            bytes_read, err := conn.Read(c.body[0:])
            if err != nil {
                log.Println("Error receiving data.")
                log.Println(err)
            } else {
                log.Printf("Received input of length %d from %v\n", bytes_read, conn.RemoteAddr().String())
                parser := new(CmdParser)
                parser.Client = &c
                err := parser.Parse(c.body[0:bytes_read])
                if err != nil {
                    log.Printf("Command parse error: %v", err.Error())
                } else {
                    server.Dispatch(parser)
                }
            }
        }(conn)
    }
}

