package pirc

import (
    "strings"
    "fmt"
)

type Channel struct {
    Name string
    Topic string
    Type string
    Mode string
    Users []*User
    UserCount int
    Operators []*User
}

func (c *Channel) UserJoin(u *User) {
    // Check if user already in channel
    for _, cu := range c.Users {
        if cu == u {
            return
        }
    }
    c.Users = append(c.Users[0:], u)
    c.UserCount += 1
    u.JoinChannel(c)
}

func (c *Channel) UserPart(u *User) {
    // Check if user is in channel and find the index in the slice
    index := 0
    found := false
    for i, cu := range c.Users[0:] {
        index = i
        if cu == u {
            found = true
            break
        }
    }

    // Delete user from channel list
    if found {
        if index == (len(c.Users) - 1) {
            c.Users = c.Users[:index]
        } else {
            c.Users = append(c.Users[:index], c.Users[index+1])
        }
        c.UserCount -= 1
    }
}

func (c *Channel) BroadcastClientCmd(sender *User, cmd string) {
    r := fmt.Sprintf(":%v!%v@%v %v\r\n", sender.Nick, sender.Username, sender.Hostname, cmd)
    for _, u := range c.Users {
        u.Conn.Write([]byte(r))
    }
}

func (c *Channel) BroadcastClientCmdNoOrig(sender *User, cmd string) {
    r := fmt.Sprintf(":%v!%v@%v %v\r\n", sender.Nick, sender.Username, sender.Hostname, cmd)
    for _, u := range c.Users {
        if u == sender { continue }
        u.Conn.Write([]byte(r))
    }
}

func (c *Channel) InChannel(u *User) bool {
    for _, cu := range c.Users {
        if cu == u {
            return true
        }
    }

    return false
}

func (c *Channel) Names(current *User) string {
    var s []string
    for _, cu := range c.Users {
        s = append(s, cu.Nick)
    }

    return strings.Join(s[0:], " ")
}


