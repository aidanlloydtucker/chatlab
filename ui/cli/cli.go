package cli

import (
	"io"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/readline.v1"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
	lg "github.com/billybobjoeaglt/chatlab/logger"
	"github.com/ttacon/chalk"
)

// The place where all the messages go
var logger *log.Logger

// Callback function for sending message
var sendMsgFunc common.SendMessageFunc

// Callback function for creating connection
var createConnFunc common.CreateConnFunc

// Username pointer
var selfUsername *string

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
	"warning":      chalk.Yellow.NewStyle().WithTextStyle(chalk.Underline).Style,
	"command":      chalk.Italic.NewStyle().Style,
}

func CLIConsole(ccChan *lg.ChanChanMessage) {
	for {
		cc := <-*ccChan
		for {
			cm, ok := <-cc
			if ok {
				switch cm.Level {
				case lg.VERBOSE:
					if lg.IsVerbose {
						logger.Println(cm.Message)
					}
				case lg.INFO:
					logger.Println(styles["notification"]("INFO: " + cm.Message))
				case lg.PRIORITY:
					logger.Println(styles["notification"]("INFO: " + cm.Message))
				case lg.WARNING:
					logger.Println(styles["warning"]("WARNING: " + cm.Message))
				case lg.ERROR:
					logger.Println(styles["error"]("ERROR: " + cm.Message))
				}
			} else {
				break
			}
		}
	}
}

// Startup function
func StartCLI() {
	// Gets username pointer from config
	selfUsername = &config.GetConfig().Username

	/*completer := readline.NewPrefixCompleter(
		readline.PcItem("/connect"),
		readline.PcItem("/connect"),
	)*/

	// Creates a readline command
	rl, err := readline.NewEx(&readline.Config{
		UniqueEditLine:  true,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		HistoryFile:     filepath.Join(common.ProgramDir, "cli-history"),
		//AutoComplete:    completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// Assign logger variable to a new logger tied to STDOUT
	logger = log.New(rl.Stdout(), "", 0)

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
		saveLine := lineHandler(line)
		if saveLine {
			rl.SaveHistory(line)
		}
	}
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
			logger.Println(styles["error"]("Error: Unknown Command"))
		}
	} else {
		if sendMsgFunc != nil && chatMap[currentChat] != nil {
			msg := common.NewMessage()
			msg.Username = *selfUsername
			msg.Message = line
			msg.ToUsers = chatMap[currentChat]
			msg.ChatName = currentChat
			AddMessage(*msg)
			sendMsgFunc(*msg)
		} else {
			logger.Println(styles["error"]("Error: No Users Connected"))
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

// Quits CLI
func QuitCLI() {
	logger.Println(styles["notification"]("Quitting"))
}

// Formats and adds message to logger
func AddMessage(msg common.Message) {
	// Message has to be to you
	if !msg.Decrypted || msg.Err != nil {
		return
	}
	// If it is a group, make it look different
	var strMsg string
	if len(msg.ToUsers) > 1 {
		strMsg += styles["group"]("(" + msg.ChatName + ") ")
	}
	strMsg += styles["username"](msg.Username+":") + " " + msg.Message
	logger.Println(strMsg)
}

// Adds a command to logger
func AddCommand(msg common.Message) {
	logger.Println(styles["command"](styles["username"](msg.Username+":") + " " + msg.Message))
}

// Removes user from usermap
func RemoveUser(user string) {
	// Go through all keys in chatMap
	for i := range chatMap {
		// For each key, see if the user is in it
		for j, val := range chatMap[i] {
			// If the user is in it, remove that user from the array
			if val == user {
				chatMap[i] = chatMap[i][:j+copy(chatMap[i][j:], chatMap[i][j+1:])]
				break
			}
		}
		// Check for empty arrays and delete them
		if len(chatMap[i]) == 0 {
			delete(chatMap, i)
		}
	}
	if user == currentChat {
		currentChat = ""
	}
	logger.Println(styles["notification"]("Removed User: " + user))
}

// Adds group to chatMap
func AddGroup(groupName string, users []string) {
	_, ok := chatMap[groupName]
	if ok && reflect.DeepEqual(chatMap[groupName], users) {
		return
	}
	if ok {
		logger.Println(styles["notification"]("Updated Group: '" + groupName + "' with the users: " + strings.Join(users, ", ")))
	} else {
		logger.Println(styles["notification"]("New Group: '" + groupName + "' with the users: " + strings.Join(users, ", ")))
	}
	chatMap[groupName] = users
	if currentChat == "" {
		currentChat = groupName
	}
}

// Adds user to chatMap
func AddUser(user string) {
	chatMap[user] = []string{user}
	if currentChat == "" {
		currentChat = user
	}
	logger.Println(styles["notification"]("New User: " + user))
}

func GetLogger() *log.Logger {
	return logger
}
