package chat

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
	"github.com/billybobjoeaglt/chatlab/crypt"
	"github.com/billybobjoeaglt/chatlab/logger"
	"github.com/billybobjoeaglt/chatlab/ui"
)

var outputChannel = make(chan chan string, 5)
var msgChan = make(chan common.Message, 5)
var peers []Peer
var peersLock = &sync.Mutex{}
var messagesReceivedAlready = make(map[string]bool)
var messagesReceivedAlreadyLock = &sync.Mutex{}

type Peer struct {
	conn     net.Conn
	username string
}

func GetOutputChannel() chan chan string {
	return outputChannel
}

func GetMessageChannel() chan common.Message {
	return msgChan
}

func CreateConnection(ip string) {
	go func() {
		conn, err := net.Dial("tcp", ip)
		if err == nil {
			handleConn(conn)
		} else {
			panic(err)
		}
	}()
}
func BroadcastMessage(msg common.Message) {
	encrypted, err := crypt.Encrypt(msg.Message, msg.ToUsers)
	if err != nil {
		panic(err)
	}
	broadcastEncryptedMessage(encrypted)
}
func broadcastEncryptedMessage(encrypted string) {
	tmpCopy := peers
	for i := range tmpCopy {
		if logger.Verbose {
			fmt.Println("Sending to " + tmpCopy[i].username)
		}
		tmpCopy[i].conn.Write([]byte(encrypted + "\n"))
	}
}
func onMessageReceived(message string, peerFrom Peer) {
	messagesReceivedAlreadyLock.Lock()
	_, found := messagesReceivedAlready[message]
	if found {
		if logger.Verbose {
			fmt.Println("Lol wait. " + peerFrom.username + " sent us something we already has. Ignoring...")
		}
		messagesReceivedAlreadyLock.Unlock()
		return
	}
	messagesReceivedAlready[message] = true
	messagesReceivedAlreadyLock.Unlock()
	//messageChannel := make(chan string, 100)
	//outputChannel <- messageChannel
	go func() {
		//defer close(messageChannel)
		processMessage(message, msgChan, peerFrom)
	}()
}
func processMessage(message string, messageChannel chan common.Message, peerFrom Peer) {
	msg := common.NewMessage()

	defer func() { messageChannel <- *msg }()

	msg.Username = peerFrom.username
	md, err := crypt.Decrypt(message)
	if err != nil {
		msg.Decrypted = false
		msg.Err = err
		return
	}
	if md.SignedBy != nil && md.SignedBy.Entity != nil && md.SignedBy.Entity.Identities != nil {
		for k := range md.SignedBy.Entity.Identities {
			/*fmt.Println("Name: " + md.SignedBy.Entity.Identities[k].UserId.Name)
			fmt.Println("Email: " + md.SignedBy.Entity.Identities[k].UserId.Email)
			fmt.Println("Comment: " + md.SignedBy.Entity.Identities[k].UserId.Comment)
			fmt.Println("Creation Time: " + md.SignedBy.Entity.Identities[k].SelfSignature.CreationTime.Format(time.UnixDate) + "\n")
			*/

			msg.Fullname = md.SignedBy.Entity.Identities[k].UserId.Name
			break
		}
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		msg.Err = err
		return
	}

	msg.Message = string(bytes)
}

func handleConn(conn net.Conn) {
	if logger.Verbose {
		fmt.Println("CONNECTION BABE. Sending our name")
	}
	conn.Write([]byte(config.GetConfig().Username + "\n"))
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	username = strings.TrimSpace(username)
	if logger.Verbose {
		fmt.Println("Received username: " + username)
	}
	//here make sure that username is valid
	peer := Peer{conn: conn, username: username}
	peersLock.Lock()
	if peerWithName(peer.username) == -1 {
		peers = append(peers, peer)
		ui.AddUser(peer.username)
		peersLock.Unlock()
		go peerListen(peer)
	} else {
		peersLock.Unlock()
		peer.conn.Close()
		if logger.Verbose {
			fmt.Println("Sadly we are already connected to " + peer.username + ". Disconnecting")
		}
	}
}
func onConnClose(peer Peer) {
	//remove from list of peers, but idk how to do that in go =(
	if logger.Verbose {
		fmt.Println("Disconnected from " + peer.username)
	}
	ui.RemoveUser(peer.username)
	peersLock.Lock()
	index := peerWithName(peer.username)
	if index == -1 {
		peersLock.Unlock()
		if logger.Verbose {
			fmt.Println("lol what")
		}
		return
	}
	peers = append(peers[:index], peers[index+1:]...)
	peersLock.Unlock()
}
func peerListen(peer Peer) {
	defer peer.conn.Close()
	defer onConnClose(peer)
	conn := peer.conn
	username := peer.username
	if logger.Verbose {
		fmt.Println("Beginning to listen to " + username)
	}
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if logger.Verbose {
				fmt.Println(err.Error())
			}
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
func Listen(port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
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
