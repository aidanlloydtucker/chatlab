package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listen()
}

func listen() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go func() {
			fmt.Println("CONNECTION BABE")
			status, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
    			fmt.Println(err)
    		}
			fmt.Printf(status)
		}()
	}
}
