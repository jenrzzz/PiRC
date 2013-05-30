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
    i := 0
    found := false
    for j, cu := range c.Users {
        i = j
        if cu == u {
            found = true
            break
        }
    }

    if found {
        u1 := c.Users[:i]
        u2 := c.Users[i+1:]
        c.Users = make([]*User, 2 * (len(u1) + len(u2)))
        copy(c.Users, u1)
        copy(c.Users[len(u1):], u2)
        c.UserCount -= 1
    }
}

func (c *Channel) BroadcastClientCmd(u *User, cmd string) {
    r := fmt.Sprintf(":%v!%v@%v %v\r\n", u.Nick, u.Username, u.Hostname, cmd)
    for _, u := range c.Users {
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
    if current != nil {
        if len(c.Users) == 0 {
            return "@" + current.Nick
        }
        s = make([]string, c.UserCount + 1)
    } else {
        s = make([]string, c.UserCount)
    }

    for _, u := range c.Users {
        if u != nil {
            s = append(s, u.Nick)
        }
    }

    if current != nil {
        s[c.UserCount] = current.Nick
        return strings.Join(s[0:c.UserCount+1], " ")
    } else {
        return strings.Join(s[0:c.UserCount], " ")
    }

    return ""
}


