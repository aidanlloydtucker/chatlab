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
	History []ChatMessage
	Name    string
	Users   []string
}

// Chat Message that extends the main message to see if user has read the message
type ChatMessage struct {
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

// To show logs before the CUI has been initalized
var chatTextInit string

func StartCUI() {
	// Add the log 'chat' to the main list
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

	chatText = termui.NewPar(chatTextInit)
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

	// Catch the keyboard
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
	if currentChat != "logs" {
		chatText.Text = ""
	}

	uiMade = true
	reloadChatList()
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

	strMsg := formatMsg(msg)

	newMsg := ChatMessage{Message: &msg, Read: msg.ChatName == currentChat}

	// Add this message to the history
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
	// If there are no users left and the user removed is the selected user,
	// nullify the currentChat and empty chat text.
	// If there are users and the user removed is the selected user,
	// go back a user and change the chat text to show that user
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
	reloadChatList()
}

// Adds group to chatMap
func AddGroup(groupName string, users []string) {
	_, ok := chatMap[groupName]
	if ok && reflect.DeepEqual(chatMap[groupName].Users, users) {
		return
	}
	// If chat already exists, notify the user that that chat has been updated
	if ok {
		orgMsg := common.NewMessage()
		orgMsg.Username = *selfUsername
		orgMsg.Message = "[INFO: Updated Group: '" + groupName + "' with the users: " + strings.Join(users, ", ") + "](fg-white,fg-underline)"
		msg := ChatMessage{Message: orgMsg, Read: false}

		tmp := chatMap[groupName]
		tmp.History = append(chatMap[groupName].History, msg)
		chatMap[msg.ChatName] = tmp
	} else {
		addToChatMap(groupName, Chat{Users: users, Name: groupName})
	}
	// If currentChat is null or it's logs, and logs is the only one besides
	// the chat just created, change the current chat
	if currentChat == "" || currentChat == "logs" && len(chatMapKeys) <= 2 {
		currentChat = groupName
	}
	reloadChatList()
}

// Adds user to chatMap
func AddUser(user string) {
	addToChatMap(user, Chat{Users: []string{user}, Name: user})

	// If currentChat is null or it's logs, and logs is the only one besides
	// the chat just created, change the current chat
	if currentChat == "" || (currentChat == "logs" && len(chatMapKeys) <= 2) {
		currentChat = user
	}
	reloadChatList()
}

// Reloads the list of chats
func reloadChatList() {
	if !uiMade {
		return
	}

	var listItems []string

	// For each chat
	for _, key := range chatMapKeys {
		var str string = "[" + key

		// If the chat is a group, add a marker
		if len(chatMap[key].Users) > 1 {
			str += " (group)"
		}

		// Count how many unread messges the chat has
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

		// Mark if there are any unread messages
		if unread > 0 {
			str += " (" + strconv.Itoa(unread) + ")"
		}

		str += "]("

		// If the chat is selected, change the styling
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
