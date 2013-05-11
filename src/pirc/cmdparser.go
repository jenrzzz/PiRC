package pirc

import (
    "bytes"
    "strings"
    "log"
    "fmt"
)

type IrcCmd struct {
    Cmd string
    Args []string
    Sender string
}

type CmdParser struct {
    Client *IrcConn
    commands []*IrcCmd
    curr int
}

func (parser *CmdParser) Next() *IrcCmd {
    if parser.curr < len(parser.commands) {
        c := parser.commands[parser.curr]
        parser.curr++
        return c
    }
    return nil
}

func (parser *CmdParser) Parse(buf []byte) error {
    cmd_list := bytes.Split(buf[0:], []byte("\n"))

    // Strip whitespace in case \r\n is sent
    var cmds []string
    if len(cmd_list[len(cmd_list)-1]) == 0 {
        cmds = make([]string, len(cmd_list)-1) // the last array from Split is empty
        for i := range cmd_list[0:len(cmd_list)-1] {
            cmds[i] = strings.TrimSpace(string(cmd_list[i]))
        }
    } else {
        // If somehow the last command did not end with a CRLF
        cmds = make([]string, len(cmd_list))
        for i := range cmd_list {
            cmds[i] = strings.TrimSpace(string(cmd_list[i]))
        }
    }

    parser.commands = make([]*IrcCmd, len(cmds))
    var err error = nil

    for i, c := range cmds {
        log.Printf("Got command %v", c)

        // Check if command includes the (optional) name of the sending server
        // Not implemented; just for RFC-1459 compliance
        var sender string
        var stripped_cmd string
        if c[0] == ':' {
            sender_end := strings.Index(c, " ")

            if sender_end == -1 {
                err = fmt.Errorf("Unable to parse command %v", c)
            } else {
                sender = c[1:sender_end]
                stripped_cmd = c[sender_end+1:]
            }
        } else {
            stripped_cmd = c
            sender = ""
        }

        cmd_split := strings.Split(stripped_cmd, " ")

        irc_cmd := IrcCmd {
            Cmd: strings.ToUpper(cmd_split[0]),
            Args: cmd_split[1:],
            Sender: sender,
        }

        parser.commands[i] = &irc_cmd
        parser.curr = 0

    }
    return err
}

