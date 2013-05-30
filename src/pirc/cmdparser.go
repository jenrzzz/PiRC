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
    Raw string
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

        split_cmds := strings.Split(stripped_cmd, " ")
        parsed_cmds := make([]string, 32)[0:len(split_cmds)]

        // Check if the last command starts with a colon (allows spaces), and
        // join those args together
        cmd_count := 0
        for j, arg := range split_cmds {
            cmd_count++
            if arg[0] == ':' {
                last_arg := strings.Join(split_cmds[j:], " ")
                parsed_cmds[j] = last_arg[1:]
                break
            } else {
                parsed_cmds[j] = split_cmds[j]
            }
        }
        parsed_cmds = parsed_cmds[0:cmd_count]

        irc_cmd := IrcCmd {
            Cmd: strings.ToUpper(parsed_cmds[0]),
            Args: parsed_cmds[1:],
            Sender: sender,
            Raw: c,
        }

        parser.commands[i] = &irc_cmd
        parser.curr = 0

    }
    return err
}

