package main

import (
"fmt"
"net"
"bufio"
)

func main() {
    fmt.Printf("Hello, world.\n")
    listen()
}

func listen(){
	ln, err := net.Listen("tcp", ":8080")
if err != nil {
	// handle error
}
for {
	conn, err := ln.Accept()
	if err != nil {
		// handle error
	}
	go func(){
		fmt.Printf("CONNECTION BABE")
		status,_:=bufio.NewReader(conn).ReadString('\n')
		
		fmt.Printf(status)
		}()
}
}

