package cli

import (
	"fmt"
	"io"
	"log"
	"strings"

	"gopkg.in/readline.v1"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
	"github.com/ttacon/chalk"
)

// The place where all the messages go
var logger *log.Logger

// Callback function for sending message
var sendMsgFunc common.SendMessageFunc

// Callback function for creating connection
var createConnFunc common.CreateConnFunc

// Username pointer
var username *string

// Map of chat name to array of users
var chatMap = make(map[string][]string)

// Current chat (index of chatMap)
var currentChat string

// Map of chalk styles
var styles = map[string]func(string) string{
	"username":     chalk.Blue.NewStyle().WithTextStyle(chalk.Bold).Style,
	"group":        chalk.Green.NewStyle().WithTextStyle(chalk.Italic).Style,
	"notification": chalk.White.NewStyle().WithTextStyle(chalk.Underline).Style,
	"error":        chalk.Red.NewStyle().WithTextStyle(chalk.Underline).Style,
	"command":      chalk.Italic.NewStyle().Style,
}

// Deprecated
func printAll(stringChanChan <-chan chan string) {
	for {
		strChan := <-stringChanChan
		for {
			str, ok := <-strChan
			if ok {
				fmt.Printf(str)
			} else {
				break
			}
		}
		fmt.Println()
	}
}

// Startup function
func StartCLI() {
	// Gets username pointer from config
	username = &config.GetConfig().Username

	// Creates a readline command
	rl, err := readline.NewEx(&readline.Config{
		UniqueEditLine:  true,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// Assign logger variable to a new logger tied to STDERR
	logger = log.New(rl.Stderr(), "", 0)

	// Set prompt
	rl.SetPrompt("> ")

	// Loop through and wait for a command
	for {
		line, err := rl.Readline()
		if err != nil {
			// If kill signal, tell main program to safe quit
			if err == readline.ErrInterrupt {
				common.Done <- true
			} else if err == io.EOF {
				common.Done <- true
			}
			break
		}
		// Handle line
		lineHandler(line)
	}
}

// Handles line
// Checks if it is a command (starts with "/") or a message
func lineHandler(line string) {
	line = strings.TrimSpace(line)

	// Must not be blank
	if line == "" {
		return
	}

	// If line is a command
	if strings.HasPrefix(line, "/") {
		// Create a message for that command
		msg := common.NewMessage()
		msg.Username = *username
		msg.Message = line
		AddCommand(*msg)

		noMatches := true
		// Checks if there is a regex match
		for _, cmd := range commandArr {
			if cmd.regex.MatchString(line) {
				// If there is a match, run a callback
				cmd.callback(line, cmd.regex.FindStringSubmatch(line)[1:])
				noMatches = false
			}
		}
		if noMatches {
			logger.Println(styles["error"]("Error: Unknown Command"))
		}
	} else {
		if sendMsgFunc != nil && chatMap[currentChat] != nil {
			msg := common.NewMessage()
			msg.Username = *username
			msg.Message = line
			msg.ToUsers = chatMap[currentChat]
			msg.ChatName = currentChat
			AddMessage(*msg)
			sendMsgFunc(*msg)
		} else {
			logger.Println(styles["error"]("Error: No Users Connected"))
		}
	}
}

func SetSendMessage(f common.SendMessageFunc) {
	sendMsgFunc = f
}

func SetCreateConn(f common.CreateConnFunc) {
	createConnFunc = f
}

func QuitCLI() {
	logger.Println(styles["notification"]("Quitting"))
}

func AddMessage(msg common.Message) {
	if !msg.Decrypted || msg.Err != nil {
		return
	}
	var strMsg string
	if len(msg.ToUsers) > 1 {
		strMsg += styles["group"]("(" + msg.ChatName + ") ")
	}
	strMsg += styles["username"](msg.Username+":") + " " + msg.Message
	logger.Println(strMsg)
}

func AddCommand(msg common.Message) {
	logger.Println(styles["command"](styles["username"](msg.Username+":") + " " + msg.Message))
}

func RemoveUser(user string) {
	for i := range chatMap {
		for j, val := range chatMap[i] {
			if val == user {
				chatMap[i] = chatMap[i][:j+copy(chatMap[i][j:], chatMap[i][j+1:])]
				break
			}
		}
	}
	if user == currentChat {
		currentChat = ""
	}
	logger.Println(styles["notification"]("Removed User: " + user))
}

func AddGroup(groupName string, users []string) {
	if chatMap[groupName] != nil {
		logger.Println(styles["notification"]("Updated Group: '" + groupName + "' with the users: " + strings.Join(users, ", ")))
	} else {
		logger.Println(styles["notification"]("New Group: '" + groupName + "' with the users: " + strings.Join(users, ", ")))
	}
	chatMap[groupName] = users
	if currentChat == "" {
		currentChat = groupName
	}
}

func AddUser(user string) {
	chatMap[user] = []string{user}
	if currentChat == "" {
		currentChat = user
	}
	logger.Println(styles["notification"]("New User: " + user))
}
