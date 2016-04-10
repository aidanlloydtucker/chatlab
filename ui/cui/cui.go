package cui

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
	"github.com/billybobjoeaglt/chatlab/utils"
	"github.com/gizak/termui"
)

// Callback function for sending message
var sendMsgFunc common.SendMessageFunc

// Callback function for creating connection
var createConnFunc common.CreateConnFunc

// Username pointer
var selfUsername *string

// Chat struct for chatMap
type Chat struct {
	History []Message
	Name    string
	Users   []string
}

type Message struct {
	*common.Message
	Read bool
}

// Map of chat name to Chat
var chatMap = make(map[string]Chat)

// List of keys for chatMap
var chatMapKeys []string

// Current chat (index of chatMap)
var currentChat string

// chatText
var chatText *termui.Par

// chatList
var chatList *termui.List

// check if us has been made
var uiMade bool

func StartCUI() {
	addToChatMap("logs", Chat{Users: []string{"logs"}, Name: "logs"})
	currentChat = "logs"
	selfUsername = &config.GetConfig().Username

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	chatInput := termui.NewPar("")
	chatInput.Border = true
	chatInput.BorderLabel = "Message"
	chatInput.Height = 3

	chatText = termui.NewPar("")
	chatText.BorderBottom = false
	chatText.BorderLeft = false
	chatText.BorderRight = false
	chatText.BorderTop = true
	chatText.Height = termui.TermHeight() - chatInput.Height - 1
	chatText.BorderFg = termui.ColorYellow

	chatList = termui.NewList()
	chatList.BorderLabel = "Chats"
	chatList.Height = termui.TermHeight() - chatInput.Height - 1
	chatList.BorderFg = termui.ColorDefault

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(4, 0, chatList),
			termui.NewCol(8, 0, chatText),
		),
		termui.NewRow(
			termui.NewCol(12, 0, chatInput),
		),
	)

	// calculate layout
	termui.Body.Align()

	termui.Render(termui.Body)

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		common.Done <- true
	})
	termui.Handle("/sys/kbd", func(ev termui.Event) {
		key := ev.Data.(termui.EvtKbd).KeyStr
		switch key {
		case "C-8", "<delete>":
			if len(chatInput.Text)-1 >= 0 {
				chatInput.Text = chatInput.Text[:len(chatInput.Text)-1]
			}
		case "<space>":
			chatInput.Text += " "
		case "<enter>":
			lineHandler(chatInput.Text)
			chatInput.Text = ""
		case "<left>":
			goLeft()
		case "<right>":
			goRight()
		case "<down>":
			goDown()
		case "<up>":
			goUp()
		default:
			chatInput.Text += key
		}
		termui.Render(termui.Body)
	})
	reloadChatList()
	uiMade = true
	termui.Loop()
}

// Handles line
// Checks if it is a command (starts with "/") or a message
func lineHandler(line string) bool {
	line = strings.TrimSpace(line)

	// Must not be blank
	if line == "" {
		return false
	}

	// If line is a command
	if strings.HasPrefix(line, "/") {
		// Create a message for that command
		msg := common.NewMessage()
		msg.Username = *selfUsername
		msg.Message = line
		AddCommand(*msg)

		noMatches := true
		// Checks if there is a regex match
		for _, cmd := range commandArr {
			if cmd.Regex.MatchString(line) {
				// If there is a match, run a callback
				cmd.Callback(line, cmd.Regex.FindStringSubmatch(line)[1:])
				noMatches = false
			}
		}
		if noMatches {
			printError("Unknown Command")
		}
	} else {
		if _, ok := chatMap[currentChat]; sendMsgFunc != nil && ok && currentChat != "logs" {
			msg := common.NewMessage()
			msg.Username = *selfUsername
			msg.Message = line
			msg.ToUsers = chatMap[currentChat].Users
			msg.ChatName = currentChat
			AddMessage(*msg)
			sendMsgFunc(*msg)
		} else {
			printError("No Users Connected")
			return false
		}
	}
	return true
}

// Sets the chat function for sending message
func SetSendMessage(f common.SendMessageFunc) {
	sendMsgFunc = f
}

// Sets the chat function for creating connection
func SetCreateConn(f common.CreateConnFunc) {
	createConnFunc = f
}

// Quits CUI
func QuitCUI() {
	printWarning("Quitting")
	termui.StopLoop()
	termui.Close()
}

// Formats and adds message to logger
func AddMessage(msg common.Message) {
	// Message has to be to you
	if !msg.Decrypted || msg.Err != nil {
		return
	}
	// If it is a group, make it look different
	strMsg := formatMsg(msg)
	newMsg := Message{Message: &msg, Read: msg.ChatName == currentChat}
	tmp := chatMap[msg.ChatName]
	tmp.History = append(chatMap[msg.ChatName].History, newMsg)
	chatMap[msg.ChatName] = tmp
	if msg.ChatName == currentChat {
		printLn(strMsg)
	} else {
		reloadChatList()
	}
}

// Removes user from usermap
func RemoveUser(user string) {
	// Go through all keys in chatMap
	for i, key := range chatMapKeys {
		// For each key, see if the user is in it
		for j, val := range chatMap[key].Users {
			// If the user is in it, remove that user from the array
			if val == user {
				tmp := chatMap[key]
				tmp.Users = append(chatMap[key].Users[:j], chatMap[key].Users[j+1:]...)
				chatMap[key] = tmp
				break
			}
		}
		// Check for empty arrays and delete them
		if len(chatMap[key].Users) == 0 {
			delete(chatMap, key)
			chatMapKeys = append(chatMapKeys[:i], chatMapKeys[i+1:]...)
		}
	}
	if user == currentChat && len(chatMapKeys) == 0 {
		currentChat = ""
		chatText.Text = ""
	} else if user == currentChat {
		currentChat = chatMapKeys[len(chatMapKeys)-1]
		var str string
		for _, msg := range chatMap[currentChat].History {
			str += formatMsg(*msg.Message) + "\n"
		}
		chatText.Text = str
	}
	if uiMade {
		reloadChatList()
	}
}

// Adds group to chatMap
func AddGroup(groupName string, users []string) {
	_, ok := chatMap[groupName]
	if ok && reflect.DeepEqual(chatMap[groupName].Users, users) {
		return
	}
	if ok {
		orgMsg := common.NewMessage()
		orgMsg.Username = *selfUsername
		orgMsg.Message = "[INFO: Updated Group: '" + groupName + "' with the users: " + strings.Join(users, ", ") + "](fg-white,fg-underline)"
		msg := Message{Message: orgMsg, Read: false}
		tmp := chatMap[groupName]
		tmp.History = append(chatMap[groupName].History, msg)
		chatMap[msg.ChatName] = tmp
	} else {
		addToChatMap(groupName, Chat{Users: users, Name: groupName})
	}
	if currentChat == "" || currentChat == "logs" && len(chatMapKeys) < 2 {
		currentChat = groupName
	}
	if uiMade {
		reloadChatList()
	}
}

// Adds user to chatMap
func AddUser(user string) {
	addToChatMap(user, Chat{Users: []string{user}, Name: user})
	if currentChat == "" || (currentChat == "logs" && len(chatMapKeys) <= 2) {
		currentChat = user
	}
	if uiMade {
		reloadChatList()
	}
}

func reloadChatList() {
	var listItems []string
	for _, key := range chatMapKeys {
		var str string
		str = "[" + key
		if len(chatMap[key].Users) > 1 {
			str += " (group)"
		}
		var unread int
		for i, msg := range chatMap[key].History {
			if !msg.Read {
				if key == currentChat {
					chatMap[key].History[i].Read = true
				} else {
					unread++
				}
			}
		}
		if unread > 0 {
			str += " (" + strconv.Itoa(unread) + ")"
		}
		str += "]("
		if key == currentChat {
			str += "fg-white,bg-green"
		}
		str += ")"

		listItems = append(listItems, str)
	}
	chatList.Items = listItems
	reloadScreen()
}

func reloadScreen() {
	termui.Render(termui.Body)
}

func addToChatMap(key string, chat Chat) {
	chatMap[key] = chat
	if !utils.ElExistsStr(chatMapKeys, key) {
		chatMapKeys = append(chatMapKeys, key)
	}
}
