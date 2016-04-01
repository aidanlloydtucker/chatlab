package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var outputChannel = make(chan chan string, 5)
var peers []Peer
var messagesReceivedAlready = make(map[string]bool)
var messagesReceivedAlreadyLock = &sync.Mutex{}
var peerChannel chan Peer

type Peer struct {
	conn     net.Conn
	username string
}
func createConnection(ip string){
	go func(){

		conn,err := net.Dial("tcp",ip)
		if err==nil{
			handleConn(conn)
		}else{
			panic(err)
		}
	}()
}
func broadcastMessage(message string){
	encrypted,err:=encrypt(message,[]string{"slaidan_lt","leijurv"})
	if err!=nil{
		panic(err)
	}
	broadcastEncryptedMessage(encrypted)
}
func broadcastEncryptedMessage(encrypted string){
	
	for i:=range peers {
		fmt.Println("Sending to "+peers[i].username)
		peers[i].conn.Write([]byte(encrypted+"\n"))
	}
}
func onMessageReceived(message string, peerFrom Peer) {
	messagesReceivedAlreadyLock.Lock()
	_, found := messagesReceivedAlready[message]
	if found {
		fmt.Println("Lol wait. " + peerFrom.username + " sent us something we already has. Ignoring...")
		messagesReceivedAlreadyLock.Unlock()
		return
	}
	messagesReceivedAlready[message] = true
	messagesReceivedAlreadyLock.Unlock()
	messageChannel := make(chan string, 100)
	outputChannel <- messageChannel
	go func() {
		defer close(messageChannel)
		processMessage(message, messageChannel, peerFrom)
	}()
}
func processMessage(message string, messageChannel chan string, peerFrom Peer) {
	messageChannel <- "Hey, a message from " + peerFrom.username + ". "
	messageChannel <- "Beginning decryption. "
	msg,err:=decrypt(message)
	if err!=nil{
		messageChannel<-"Unable to decrypt =("
		messageChannel<-err.Error()
		return
	}
	messageChannel <- "Done decrypting. "
	messageChannel <- "Here's the message: "
	messageChannel <- msg
}

func handleConn(conn net.Conn) {
	fmt.Println("CONNECTION BABE. Sending our name")
	conn.Write([]byte(config.Username + "\n"))
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	username = strings.TrimSpace(username)
	fmt.Println("Received username: " + username)
	//here make sure that username is valid
	peerObj := Peer{conn: conn, username: username}
	peerChannel <- peerObj
}
func onConnClose(peer Peer) {
	//remove from list of peers, but idk how to do that in go =(
	fmt.Println("Disconnected from " + peer.username)
}
func peerListen(peer Peer) {
	defer peer.conn.Close()
	defer onConnClose(peer)
	conn := peer.conn
	username := peer.username
	fmt.Println("Beginning to listen to " + username)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		message = strings.TrimSpace(message)
		onMessageReceived(message, peer)
	}
}
func peerWithName(name string) int {
	for i := 0; i < len(peers); i++ {
		if peers[i].username == name {
			return i
		}
	}
	return -1
}
func listen() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	peerChannel = make(chan Peer)
	defer close(peerChannel)
	go func() {
		for {
			peer, ok := <-peerChannel
			if ok {
				if peerWithName(peer.username) == -1 {
					peers = append(peers, peer)
					go peerListen(peer)
				} else {
					peer.conn.Close()
					fmt.Println("Sadly we are already connected to " + peer.username + ". Disconnecting")
				}
			} else {
				fmt.Println("Peers over")
				return
			}
		}
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}
