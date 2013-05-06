package main

import (
    "fmt"
    "pirc"
)

func main() {
    fmt.Println("pIRC starting on port 6667...")
    pirc.RunServer(":6667")
}
