package pirc

import (
    "fmt"
    "log"
    "net"
)

func RunServer() {
    port := ":6667"
    ln, err := net.Listen("tcp", port)
    if err != nil {
        log.Println(err)
        log.Panicf("Unable to start a server on %v!", port)
    }

    // Main connection loop
    for {
        log.Println("Waiting for connection")
        conn, err := ln.Accept()
        if err != nil {
            log.Println("Error opening connection.")
            log.Println(err)
            continue
        } else {
            log.Println("Connection opened")
        }

        go func(c net.Conn) {
            var msg [1024]byte
            chars, err := c.Read(msg[0:])
            if err != nil {
                log.Println("Error receiving data.")
                log.Println(err)
            } else {
                fmt.Printf("Received input from %v\n", c.RemoteAddr().String())
                fmt.Println(string(msg[:chars]))
            }
        }(conn)
    }
}
