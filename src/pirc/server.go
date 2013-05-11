package pirc

import (
    "log"
    "net"
)

type IrcConn struct {
    rwc net.Conn
    remoteAddr string
    server *Server
    body []byte
}

func (c *IrcConn) Read(b []byte) (n int, err error) {
    return c.rwc.Read(b)
}

func (c *IrcConn) Write(b []byte) (n int, err error) {
    return c.rwc.Write(b)
}

type Server struct {
    Hostname string
    Users map[string] *User
    Connections map[*IrcConn] *User
    // channels map[string] *Channel
}

// Check if command is valid and dispatch, then write response
func (s *Server) Dispatch(cmds *CmdParser) {
    log.Printf("In dispatch with commands %v", cmds.commands)
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
    server.Users = make(map[string] *User)
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
            c := IrcConn{conn, conn.RemoteAddr().String(), server, make([]byte, 2048)}
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

