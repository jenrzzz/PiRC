package pirc

import (
    "regexp"
    "log"
)

var CmdDispatcher = map[string] func(*IrcCmd, *IrcConn) ServerResponse {
    "NICK": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        // Make sure nick was provided
        if len(cmd.Args) < 1 {
            return ERR["NEEDMOREPARAMS"]
        }

        // Check if nick is valid
        nick := cmd.Args[0]
        if valid, err := regexp.MatchString("[A-Za-z][A-Za-z0-9\\[\\]\\`\\^\\{\\}\\-]*", nick); err != nil || !valid {
            log.Printf("Invalid nickname %v", nick)
            return ERR["ERRONEOUSNICKNAME"]
        } 

        // Check whether this is a new connection or user is changing their nick
        if conn_user, exists := server.Connections[conn]; exists {
            return server.ChangeNick(conn_user, nick)
        }

        if server.UsersByNick[nick] != nil {
            return ERR["NICKNAMEINUSE"]
        }

        u, _ := CreateUser(nick, conn)
        server.Connections[conn] = u

        log.Printf("Created user %v", *u)
        return nil
    },

    "USER": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        if len(cmd.Args) < 4 {
            return ERR["NEEDMOREPARAMS"]
        }

        username := cmd.Args[0]
        // don't care about self-reported hostname (2nd arg)
        servername := cmd.Args[2]
        realname := cmd.Args[3]

        if server.FindUserByName(username) != nil {
            return ERR["ALREADYREGISTRED"]
        }

        u := server.FindUserByConn(conn)
        u.Username = username
        u.Servername = servername
        u.Realname = realname

        server.RegisterUser(u)
        log.Printf("Registered user %v", *u)
        return nil
    },
}
