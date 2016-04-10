package chat

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/crypt"
	"github.com/billybobjoeaglt/chatlab/logger"
	"github.com/billybobjoeaglt/chatlab/ui"
)

// XXX: DEPRICATED
var outputChannel = make(chan chan string, 5)

// Channel for adding messages
var msgChan = make(chan common.Message, 5)

// List of all connected peers
var peers []Peer
var peersLock = &sync.Mutex{}

// Map of recieved messages
// TODO: Make the key a hash of the encrypted message
var messagesReceivedAlready = make(map[string]bool)
var messagesReceivedAlreadyLock = &sync.Mutex{}

// The node of the program
var SelfNode = Node{}

// The struct that is sent to peers
type EncyptedMessage struct {
	EncyptedMessage string
}

// Defines what a connection is
type Node struct {
	Username string
	IsRelay  bool
	Port     string
}

// Defines who the connection is and how to understand them
type Peer struct {
	Conn     net.Conn
	Username string
	Encoder  *gob.Encoder
	Decoder  *gob.Decoder
	Node     Node
}

// XXX: DEPRICATED
func GetOutputChannel() chan chan string {
	return outputChannel
}

// Returns msgChan
func GetMessageChannel() chan common.Message {
	return msgChan
}

// Adds a connection to the peers
func CreateConnection(ip string, silent bool) {
	// Creates a channel for logging
	cc := make(logger.ChanMessage)
	logger.ConsoleChan <- cc
	go func() {
		conn, err := net.Dial("tcp", ip)
		if err == nil {
			handleConn(conn, cc, silent)
		} else {
			cc.AddError(err, "Could not connect")
			close(cc)
		}
	}()
}

// Takes message, encrypts it, and sends it to all peers
func BroadcastMessage(msg common.Message) {
	encrypted, err := crypt.EncryptMessage(msg)
	if err != nil {
		panic(err)
	}
	encMsg := EncyptedMessage{EncyptedMessage: encrypted}
	broadcastEncryptedMessage(encMsg)
}

// Sends an encrypted message to all peers
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

// Runs when a peer sends a message.
// Checks whether the message has already been recieved.
// Sets the message to be recieved and proccesses it.
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

// Proccesses an encrypted message.
// Decrypts the message and adds it to the ui.
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

// Called on the start of a connection.
// Builds the Peer struct.
// Sends back the local node.
// Runs a listener thread if username exists.
func handleConn(conn net.Conn, cc logger.ChanMessage, silent bool) {
	defer close(cc)
	cc.AddVerbose("Received connection. Sending self data")

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// Sends local node back through connection
	err := encoder.Encode(SelfNode)
	if err != nil {
		if !silent || logger.IsVerbose {
			cc.AddError(err, "Count not encode SelfNode")
		}
		return
	}

	node := Node{}
	err = decoder.Decode(&node)
	if err != nil {
		if !silent || logger.IsVerbose {
			cc.AddError(err, "Could not decode node gob")
		}
		return
	}

	cc.AddVerbose("Received username: " + node.Username + " Relay: " + strconv.FormatBool(node.IsRelay))

	peer := Peer{Conn: conn, Username: node.Username, Encoder: encoder, Decoder: decoder, Node: node}

	peersLock.Lock()
	// If a peer with the same username does not exist
	if peerWithName(peer.Username) == -1 {
		// Add the peer to the list
		peers = append(peers, peer)
		peersLock.Unlock()

		// If the peer is not a relay, add it to the ui
		if !peer.Node.IsRelay {
			ui.AddUser(peer.Username)
		}

		// Start a process to listen to the peer
		go peerListen(peer)

		// Log that you conencted
		if !silent || logger.IsVerbose {
			cm := logger.ConsoleMessage{Level: logger.PRIORITY}
			cm.Message = "Connected to "
			if node.IsRelay {
				cm.Message += "Relay"
			} else {
				cm.Message += "Node"
			}
			cc <- cm
		}
	} else {
		peersLock.Unlock()
		// Close connection with peer
		peer.Conn.Close()

		// Log that you are already connected
		if !silent || logger.IsVerbose {
			cm := logger.ConsoleMessage{Level: logger.PRIORITY}
			cm.Message = "Already Connected to "
			if node.IsRelay {
				cm.Message += "Relay"
			} else {
				cm.Message += "Node"
			}
			cc <- cm
		}
	}
}

// When the peer (or hose) closes the connection
func onConnClose(peer Peer) {
	logger.ConsoleChan.AddVerbose("Disconnected from peer: " + peer.Username)
	// If peer was not a relay, remove from UI
	if !peer.Node.IsRelay {
		ui.RemoveUser(peer.Username)
	}

	// Remove the peer from the list of peers
	peersLock.Lock()
	index := peerWithName(peer.Username)
	if index == -1 {
		peersLock.Unlock()
		logger.ConsoleChan.AddError(nil, "Tried to disconnect from peer but it was already disconnected")
		return
	}
	peers = append(peers[:index], peers[index+1:]...)
	peersLock.Unlock()
}

// Listen to peer
func peerListen(peer Peer) {
	defer peer.Conn.Close()
	defer onConnClose(peer)
	logger.ConsoleChan.AddVerbose("Beginning to listen to " + peer.Username)

	// When the peer sends a message, decode it and send it off to be parsed
	for {
		encMsg := &EncyptedMessage{}
		err := peer.Decoder.Decode(encMsg)
		if err != nil {
			return
		}
		onMessageReceived(*encMsg, peer)
	}
}

// Search through peers list for the index of a peer by the username
func peerWithName(name string) int {
	for i := 0; i < len(peers); i++ {
		if peers[i].Username == name {
			return i
		}
	}
	return -1
}

// Format for peers to be saved to disk
type SavedPeer struct {
	IP       string
	Username string
	Key      string
	IsRelay  bool
}

// This saves all peers to disk.
// It goes through all currently connected peers and saves them to ~/.chatlab/saved-peers.gob
func SavePeers() error {
	var savedPeers []SavedPeer

	// Go through each connected peer
	peersLock.Lock()
	for _, peer := range peers {
		// Get IP Address
		tcpAddrIP := peer.Conn.RemoteAddr().(*net.TCPAddr).IP.String()

		// Create a SavedPeer with the ip, port, username, and if it's a relay
		savedPeer := SavedPeer{IP: net.JoinHostPort(tcpAddrIP, peer.Node.Port), Username: peer.Username, IsRelay: peer.Node.IsRelay}

		// If it is not a relay and a key exists, save the public key too
		if _, ok := crypt.GetKeyMap()[peer.Username]; !savedPeer.IsRelay && ok {
			pkBuf := new(bytes.Buffer)
			crypt.GetKeyMap()[peer.Username].Serialize(pkBuf)
			savedPeer.Key = pkBuf.String()
		}

		// Add peer to the list
		savedPeers = append(savedPeers, savedPeer)
	}
	peersLock.Unlock()

	// Encode list to file
	var gobBuf bytes.Buffer
	enc := gob.NewEncoder(&gobBuf)
	err := enc.Encode(savedPeers)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(common.ProgramDir, "saved-peers.gob"), gobBuf.Bytes(), 0777)
	return err
}

// This loads saved peers from disk
func LoadPeers() error {
	peerPath := filepath.Join(common.ProgramDir, "saved-peers.gob")

	// If file exists
	if _, err := os.Stat(peerPath); err == nil {
		file, err := os.Open(peerPath)
		if err != nil {
			return err
		}
		defer file.Close()

		dec := gob.NewDecoder(file)

		var savedPeers []SavedPeer
		err = dec.Decode(&savedPeers)
		if err != nil {
			return err
		}

		// Go through all the saved peers and create a connection
		for _, savedPeer := range savedPeers {
			go CreateConnection(savedPeer.IP, true)

			// If the peer had a key and was not a relay, add the public key
			if !savedPeer.IsRelay && savedPeer.Key != "" {
				go crypt.AddPublicKeyToMap(savedPeer.Username, savedPeer.Key)
			}
		}
	}

	return nil
}

// Create a listener for the app to lisen on
func Listen(port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// For each new connection, handle the connection in a new process
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		cc := make(logger.ChanMessage)
		logger.ConsoleChan <- cc
		go handleConn(conn, cc, false)
	}
}
