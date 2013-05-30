package pirc

import (
    "log"
    "net"
    "strings"
    "fmt"
    "github.com/stathat/treap"
    "time"
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
    log.Printf("Writing response %v\n", string(b))
    return c.rwc.Write(b)
}

func (c *IrcConn) WriteRaw(code int, nick string, message string) (n int, err error) {
    r := fmt.Sprintf(":%v %03d %v %v\r\n", c.server.Hostname, code, nick, message)
    return c.Write([]byte(r))
}

func (c *IrcConn) WriteResponse(r ServerResponse, cmd string) (n int, err error) {
    return c.Write([]byte(r.Response(c.server, cmd)))
}

func (c *IrcConn) WriteFormattedResponse(s *Server, code int, nick string, resp string) (n int, err error) {
    response := fmt.Sprintf(":%v %03d %v %v\r\n", s.Hostname, code, nick, resp)
    return c.Write([]byte(response))
}

func (c *IrcConn) WriteCmd(cmd string, args []string) (n int, err error) {
    var r string
    if args != nil {
        r = fmt.Sprintf(":%v %v %v\r\n", c.server.Hostname, cmd, strings.Join(args[0:], " "))
    } else {
        r = fmt.Sprintf(":%v %v\r\n", c.server.Hostname, cmd)
    }
    return c.Write([]byte(r))
}

func (c *IrcConn) WriteServerNotice(s string) (n int, err error) {
    r := fmt.Sprintf(":%v NOTICE * :%v\r\n", c.server.Hostname, s)
    return c.Write([]byte(r))
}

// Comparator for treap
func StrLess(p, q interface{}) bool {
    a, b := p.(string), q.(string)
    var min = len(b)
    if len(a) < len(b) {
        min = len(a)
    }

    var diff int
    for i := 0; i < min && diff == 0; i++ {
        diff = int(a[i]) - int(b[i])
    }

    if diff == 0 {
        diff = len(a) - len(b)
    }

    return diff < 0
}

// Das server
type Server struct {
    Hostname string
    UsersByNick map[string] *User
    UsersByUsername map[string] *User
    Connections map[*IrcConn] *User
    StartedAt time.Time
    Version string
    Channels *treap.Tree
}

func (s *Server) RegisterConnection(u *User) {
    s.Connections[u.Conn] = u
}

func (s *Server) UnregisterConnection(c *IrcConn) {
    u := s.FindUserByConn(c)
    delete(s.Connections, c)
    delete(s.UsersByNick, u.Nick)
    delete(s.UsersByUsername, u.Username)
}


func (s *Server) RegisterUser(u *User) {
    if server.FindUserByName(u.Username) != nil {
        u.Conn.WriteResponse(ERR["ALREADYREGISTRED"], u.Nick)
    }

    s.UsersByNick[u.Nick] = u
    s.UsersByUsername[u.Username] = u
    s.Connections[u.Conn] = u
    u.Registered = true

    // Write connection banner
    u.Conn.WriteResponse(RPL["WELCOME"].Format(u.FullyQualifiedName()), u.Nick)
    u.Conn.WriteResponse(RPL["YOURHOST"].Format(s.Hostname, s.Version), u.Nick)
    u.Conn.WriteResponse(RPL["CREATED"].Format(s.StartedAt), u.Nick)
    u.Conn.WriteResponse(RPL["MYINFO"].Format(s.Hostname, s.Version, "opsitnmlbvk", "iswo"), u.Nick)
    u.Conn.WriteResponse(RPL["MOTDSTART"].Format(s.Hostname), u.Nick)
    u.Conn.WriteResponse(RPL["ENDOFMOTD"], u.Nick)
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
        var response ServerResponse

        // Check if command is valid
        if cmd_func, exists := CmdDispatcher[c.Cmd]; !exists {
            response = ERR["UNKNOWNCOMMAND"].Format(c.Cmd)
        } else {
            // Check if user is registered if they aren't trying to register
            if c.Cmd != "NICK" && c.Cmd != "USER" {
                u := s.FindUserByConn(cmds.Client)
                if u == nil || !u.Registered {
                    response = ERR["NOTREGISTERED"]
                } else {
                    response = cmd_func(c, cmds.Client)
                }
            } else {
                response = cmd_func(c, cmds.Client)
            }
        }

        if response != nil {
            cmds.Client.WriteResponse(response, c.Cmd)
        }
    }
}

var server = new(Server)
func RunServer(listenaddr string) {
    server.UsersByNick = make(map[string] *User)
    server.UsersByUsername = make(map[string] *User)
    server.Connections = make(map[*IrcConn] *User)
    server.Channels = treap.NewTree(StrLess)
    server.Hostname = "localhost"
    server.Version = "piRC 0.0.1-alpha1"
    server.StartedAt = time.Now()
    ln, err := net.Listen("tcp", listenaddr)
    if err != nil {
        log.Println(err)
        log.Panicf("Unable to start a server on %v!", listenaddr)
    }

    // Connection handler loop
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

        // Connection processor loop
        // No select(2) -- it's a beautiful thing :D
        // conn.Read will block in this goroutine until data is received
        go func(conn net.Conn) {
            remote_addr := conn.RemoteAddr().String()
            remote_addr = remote_addr[0:strings.IndexAny(remote_addr, ":")]
            c := IrcConn{conn, remote_addr, server, make([]byte, 2048)}

            for {
                bytes_read, err := conn.Read(c.body[0:])
                if err != nil {
                    if err.Error() == "EOF" {
                        log.Printf("Connection from %v closed.", remote_addr)
                    } else {
                        log.Println("Error receiving data.")
                        log.Println(err)
                    }
                    server.UnregisterConnection(&c)
                    break
                } else {
                    parser := new(CmdParser)
                    parser.Client = &c
                    err := parser.Parse(c.body[0:bytes_read])
                    if err != nil {
                        log.Printf("Command parse error: %v", err.Error())
                    } else {
                        server.Dispatch(parser)
                    }
                }
                c.body = make([]byte, 2048)
            }
        }(conn)
    }
}

