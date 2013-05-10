package pirc

import (
    "bytes"
    "strings"
    "log"
)

type IrcCmd struct {
    Cmd string
    Args []string
}

type CmdParser struct {
    commands []*IrcCmd
}

func (parser *CmdParser) Parse(buf []byte) (string, *CodePair) {
    log.Printf("In: %v", buf)
    cmd_list := bytes.Split(buf[0:], []byte("\n"))
    log.Printf("Split: %v", cmd_list)

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

    for _, c := range cmds {
        log.Printf("Got command %v", c)
        cmd_split := strings.Split(c, " ")

        irc_cmd := IrcCmd {
            Cmd: cmd_split[0],
            Args: cmd_split[1:],
        }

        // Check if command is valid
        if _, ok := CmdDispatcher[irc_cmd.Cmd]; !ok {
            return irc_cmd.Cmd, &ERR.UNKNOWNCOMMAND
        }

    }
    return "", nil
}

var CmdDispatcher = make(map[string] func(*IrcCmd) ServerResponse)
