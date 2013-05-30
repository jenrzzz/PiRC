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

    "PING": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        if len(cmd.Args) < 1 {
            return ERR["NEEDMOREPARAMS"]
        }

        conn.WriteCmd("PONG", nil)
        return nil
    },

    "PONG": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        log.Printf("%v is still alive", conn.remoteAddr)
        return nil
    },

    "JOIN": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        channel := cmd.Args[0]
        u := server.FindUserByConn(conn)
        var c *Channel
        if !server.Channels.Exists(channel) {
            c = new(Channel)
            c.Name = channel
            c.Operators = []*User{u}
            c.Users = []*User{u}
            server.Channels.Insert(channel, c)
        } else {
            c = server.Channels.Get(channel).(*Channel)
        }

        c.UserJoin(u)

        if c.Topic != "" {
            conn.WriteFormattedResponse(server, RPL["TOPIC"].Code, u.Nick, RPL["TOPIC"].FormatMsg(channel, c.Topic))
        } else {
            conn.WriteFormattedResponse(server, RPL["NOTOPIC"].Code, u.Nick, RPL["NOTOPIC"].FormatMsg(channel))
        }

        conn.WriteFormattedResponse(server, RPL["NAMREPLY"].Code, u.Nick, RPL["NAMREPLY"].FormatMsg(channel, c.Names(u)))
        conn.WriteFormattedResponse(server, RPL["ENDOFNAMES"].Code, u.Nick, RPL["ENDOFNAMES"].FormatMsg(channel))

        c.BroadcastClientCmd(u, cmd.Raw)

        return nil
    },

    "PART": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        if len(cmd.Args) < 1 {
            return ERR["NEEDMOREPARAMS"]
        }

        channel := server.Channels.Get(cmd.Args[0]).(*Channel)
        u := server.FindUserByConn(conn)
        if u == nil {
            return ERR["NOSUCHCHANNEL"]
        }

        if ! channel.InChannel(u) {
            return ERR["NOTONCHANNEL"]
        }

        channel.BroadcastClientCmd(u, cmd.Raw)
        channel.UserPart(u)
        return nil
    },

    "MODE": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        if len(cmd.Args) < 1 {
            return ERR["NEEDMOREPARAMS"]
        }

        // Don't support channel modes yet
        if cmd.Args[0][0] == '#' {
            return nil
        }

        u := server.FindUserByConn(conn)
        if u != server.FindUserByNick(cmd.Args[0]) {
            return ERR["USERSDONTMATCH"]
        }

        return nil
    },
}
