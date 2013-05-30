package pirc

import (
    "fmt"
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
    Channels []*Channel
    Conn *IrcConn
}

func (u *User) FullyQualifiedName() string {
    return fmt.Sprintf("%v!%v@%v", u.Nick, u.Username, u.Hostname)
}

func (u *User) JoinChannel(c *Channel) {
    u.Channels = append(u.Channels, c)
}

func (u *User) PartChannel(c *Channel) {
    i := 0
    found := false
    for j, cc := range u.Channels {
        i = j
        if cc == c {
            found = true
            break
        }
    }

    if found {
        c1 := u.Channels[:i]
        c2 := u.Channels[i+1:]
        u.Channels = make([]*Channel, 2 * (len(c1) + len(c2)))
        copy(u.Channels, c1)
        copy(u.Channels[len(c1):], c2)
    }
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


