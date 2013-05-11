package pirc

var CmdDispatcher = map[string] func(*IrcCmd, *IrcConn) ServerResponse {
    "NICK": func(cmd *IrcCmd, conn *IrcConn) ServerResponse {
        // Check whether this is a new connection or user is changing their nick
        var u *User
        if conn_user, exists := server.Connections[conn]; exists {
            // Delete the old nick if we're changing it
            u = conn_user
            delete(server.Users, u.Nick)
        } else {
            u = new(User)
        }

        if len(cmd.Args) < 1 {
            return ERR["NEEDMOREPARAMS"]
        }

        nick := cmd.Args[0]
        if server.Users[nick] != nil {
            return ERR["NICKNAMEINUSE"]
        }

        u.Nick = nick
        server.Users[nick] = u
        server.Connections[conn] = u
        return nil
    },
}
