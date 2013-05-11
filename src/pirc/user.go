package pirc

import (
    "net"
    "log"
)

type User struct {
    Nick string
    Username string
    Realname string
    Hostname string
    Servername string
    Registered bool
    Conn *IrcConn
}

func CreateUser(nick string, c *IrcConn) (*User, error) {
    u := new(User)
    u.Nick = nick
    u.Conn = c

    // Lookup hostname
    addrs, err := net.LookupAddr(c.remoteAddr)
    if err != nil {
        log.Printf("Failed to lookup hostname for %v", c.remoteAddr)
        u.Hostname = c.remoteAddr
    } else {
        u.Hostname = addrs[0]
    }

    return u, nil
}
