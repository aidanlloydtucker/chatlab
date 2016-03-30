package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	temp:=make(chan string)
	printAll(temp)
	temp<-"lol"
	listen()
}

func handleConn(conn net.Conn) {
    fmt.Println("CONNECTION BABE")
    status, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf(status)
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
		go handleConn(conn)
	}
}

func printAll(stringChan <-chan string) {
	go func(){
		for{
			str:=<-stringChan
			fmt.Printf(str)
		}
	}()
}