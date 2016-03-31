package main

import (
	"bufio"
	"fmt"
	"strings"
	"net"
	//"golang.org/x/crypto/openpgp"
)

var outputChannel = make(chan chan string, 5)
var peers []Peer
type Peer struct {
	conn net.Conn
	username string
}
func main() {
	go printAll(outputChannel)
	listen()
}
func onMessageReceived(message string, peerFrom Peer) {
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

func handleConn(conn net.Conn, peerChannel chan Peer) {
	fmt.Println("CONNECTION BABE")
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	username=strings.TrimSpace(username)
	fmt.Println("Received username: "+username)
	peerObj:=Peer{conn:conn,username:username}
	peerChannel<-peerObj
}

func peerListen(peer Peer){
	conn:=peer.conn
	username:=peer.username
	defer conn.Close()
	fmt.Println("Beginning to listen to "+username)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err!=nil{
			return
		}
		message=strings.TrimSpace(message)
		onMessageReceived(message,peer)
	}
}

func listen() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	peerChannel := make(chan Peer)
	defer close(peerChannel)
	go func(){
		for{
			peer,ok := <-peerChannel
			if ok{
				//here check if we are already connected to the same username and if so close the connection
				peers = append(peers,peer)
				go peerListen(peer)
			}else{
				return
			}
		}
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn,peerChannel)
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
