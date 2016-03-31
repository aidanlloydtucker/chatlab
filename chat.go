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

type Peer struct {
	conn     net.Conn
	username string
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
	messageChannel <- "Beginning processsing. "
	messageChannel <- "Done processing. "
	messageChannel <- "Here's the message: "
	messageChannel <- message
}

func handleConn(conn net.Conn, peerChannel chan Peer) {
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
	peerChannel := make(chan Peer)
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
				return
			}
		}
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn, peerChannel)
	}
}
