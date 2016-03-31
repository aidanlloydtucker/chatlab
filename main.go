package main

import (
	"bufio"
	"fmt"
	"strings"
	"net"
	//"golang.org/x/crypto/openpgp"
)

var outputChannel = make(chan chan string, 5)
func main() {
	go printAll(outputChannel)
	onMessageReceived(outputChannel, "msg1")
	onMessageReceived(outputChannel, "msg2")
	onMessageReceived(outputChannel, "msg3")
	onMessageReceived(outputChannel, "msg4")
	listen()
}
func onMessageReceived(outputChannel chan chan string, message string) {
	messageChannel := make(chan string, 100)
	outputChannel <- messageChannel
	go func(){
		defer close(messageChannel)
	 	processMessage(message,messageChannel)
	}()
}
func processMessage(message string, messageChannel chan string) {
	messageChannel<-"Beginning processsing. "
	messageChannel<-"Done processing. "
	messageChannel<-"Here's the message: "
	messageChannel<-message
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("CONNECTION BABE")
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	username=strings.TrimSpace(username)
	fmt.Println("Username: "+username)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err!=nil{
			panic(err)
		}
		message=strings.TrimSpace(message)
		fmt.Println("Message from "+username+": "+message)
	}
}

func listen() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func printAll(stringChanChan <-chan chan string) {
	for {
		strChan := <-stringChanChan
		for{
			str, ok:= <-strChan
			if ok{
				fmt.Printf(str)
			}else{
				break
			}
		}
		fmt.Println()
	}
}
