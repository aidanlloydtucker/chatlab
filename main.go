package main

import (
	"bufio"
	"fmt"
	"net"
	//"golang.org/x/crypto/openpgp"
)

func main() {
	outputChannel := make(chan chan string, 5)
	printAll(outputChannel)
	onMessageReceived(outputChannel, "lol")
	listen()
}
func onMessageReceived(outputChannel chan chan string, message string) {
	messageChannel := make(chan string)
	outputChannel <- messageChannel
	go func() {
		messageChannel <- processMessage(message)
	}()
}
func processMessage(message string) string {
	return "Processed Message: " + message
}

func handleConn(conn net.Conn) {
	fmt.Println("CONNECTION BABE")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(status)
}

func listen() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func printAll(stringChanChan <-chan chan string) {
	go func() {
		for {
			strChan := <-stringChanChan
			str := <-strChan
			fmt.Println(str)
		}
	}()
}
