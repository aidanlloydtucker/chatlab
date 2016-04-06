package chat

import (
	"encoding/gob"
	"net"
	"strconv"
	"sync"

	"github.com/billybobjoeaglt/chatlab/common"
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
var SelfNode = Node{}

type EncyptedMessage struct {
	EncyptedMessage string
}

type Node struct {
	Username string
	IsRelay  bool
}

type Peer struct {
	Conn     net.Conn
	Username string
	Encoder  *gob.Encoder
	Decoder  *gob.Decoder
	Node     Node
}

func GetOutputChannel() chan chan string {
	return outputChannel
}

func GetMessageChannel() chan common.Message {
	return msgChan
}

func CreateConnection(ip string) {
	cc := make(logger.ChanMessage)
	logger.ConsoleChan <- cc
	go func() {
		conn, err := net.Dial("tcp", ip)
		if err == nil {
			handleConn(conn, cc)
		} else {
			cc.AddError(err, "Could not connect")
			close(cc)
		}
	}()
}
func BroadcastMessage(msg common.Message) {
	encrypted, err := crypt.EncryptMessage(msg)
	if err != nil {
		panic(err)
	}
	encMsg := EncyptedMessage{EncyptedMessage: encrypted}
	broadcastEncryptedMessage(encMsg)
}
func broadcastEncryptedMessage(encMsg EncyptedMessage) {
	messagesReceivedAlreadyLock.Lock()
	messagesReceivedAlready[encMsg.EncyptedMessage] = true
	messagesReceivedAlreadyLock.Unlock()
	tmpCopy := peers
	for i := range tmpCopy {
		logger.ConsoleChan.AddVerbose("Sending encrypted message to " + tmpCopy[i].Username)
		tmpCopy[i].Encoder.Encode(encMsg)
	}
}
func onMessageReceived(encMsg EncyptedMessage, peerFrom Peer) {
	messagesReceivedAlreadyLock.Lock()
	_, found := messagesReceivedAlready[encMsg.EncyptedMessage]
	if found {
		logger.ConsoleChan.AddVerbose("A peer (" + peerFrom.Username + ") sent us something we already have. Ignoring...")
		messagesReceivedAlreadyLock.Unlock()
		return
	}
	messagesReceivedAlready[encMsg.EncyptedMessage] = true
	messagesReceivedAlreadyLock.Unlock()
	//messageChannel := make(chan string, 100)
	//outputChannel <- messageChannel
	broadcastEncryptedMessage(encMsg)
	go func() {
		//defer close(messageChannel)
		processMessage(encMsg, msgChan, peerFrom)
	}()
}
func processMessage(encMsg EncyptedMessage, messageChannel chan common.Message, peerFrom Peer) {

	md, msg, err := crypt.DecryptMessage(encMsg.EncyptedMessage)
	if err != nil {
		return
	}

	defer func() { messageChannel <- *msg }()

	if len(msg.ToUsers) > 1 {
		ui.AddGroup(msg.ChatName, msg.ToUsers)
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
}

func handleConn(conn net.Conn, cc logger.ChanMessage) {
	defer close(cc)
	cc.AddVerbose("Received connection. Sending self data")

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	err := encoder.Encode(SelfNode)
	if err != nil {
		cc.AddError(err, "Count not encode SelfNode")
		return
	}

	node := Node{}
	err = decoder.Decode(&node)
	if err != nil {
		cc.AddError(err, "Could not decode node gob")
		return
	}
	cc.AddVerbose("Received username: " + node.Username + " Relay: " + strconv.FormatBool(node.IsRelay))

	//here make sure that username is valid
	peer := Peer{Conn: conn, Username: node.Username, Encoder: encoder, Decoder: decoder, Node: node}
	peersLock.Lock()
	if peerWithName(peer.Username) == -1 {
		peers = append(peers, peer)
		if !peer.Node.IsRelay {
			ui.AddUser(peer.Username)
		}
		peersLock.Unlock()
		go peerListen(peer)

		cm := logger.ConsoleMessage{Level: logger.INFO}
		cm.Message = "Connected to "
		if node.IsRelay {
			cm.Message += "Relay"
		} else {
			cm.Message += "Node"
		}
		cc <- cm
	} else {
		peersLock.Unlock()
		peer.Conn.Close()

		cm := logger.ConsoleMessage{Level: logger.INFO}
		cm.Message = "Already Connected to "
		if node.IsRelay {
			cm.Message += "Relay"
		} else {
			cm.Message += "Node"
		}
		cc <- cm
	}
}
func onConnClose(peer Peer) {
	//remove from list of peers, but idk how to do that in go =(
	logger.ConsoleChan.AddVerbose("Disconnected from peer: " + peer.Username)
	if !peer.Node.IsRelay {
		ui.RemoveUser(peer.Username)
	}
	peersLock.Lock()
	index := peerWithName(peer.Username)
	if index == -1 {
		peersLock.Unlock()
		logger.ConsoleChan.AddVerbose("Lol What? Index is -1")
		return
	}
	peers = append(peers[:index], peers[index+1:]...)
	peersLock.Unlock()
}

// XXX: For some reason logger cannot be called here without it breaking everything
func peerListen(peer Peer) {
	defer peer.Conn.Close()
	defer onConnClose(peer)
	logger.ConsoleChan.AddVerbose("Beginning to listen to " + peer.Username)
	for {
		encMsg := &EncyptedMessage{}
		err := peer.Decoder.Decode(encMsg)
		if err != nil {
			return
		}
		onMessageReceived(*encMsg, peer)
	}
}
func peerWithName(name string) int {
	for i := 0; i < len(peers); i++ {
		if peers[i].Username == name {
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
		cc := make(logger.ChanMessage)
		logger.ConsoleChan <- cc
		go handleConn(conn, cc)
	}
}
