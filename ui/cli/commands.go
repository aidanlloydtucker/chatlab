package cli

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/billybobjoeaglt/chatlab/common"
)

type CommandCallback func(string, []string)

type Command struct {
	regex    *regexp.Regexp
	command  string
	desc     string
	args     string
	example  []string
	callback CommandCallback
}

var commandArr = []Command{
	Command{
		regex:   regexp.MustCompile(`\/connect ([^ ]+) ?(.*)`),
		command: "connect",
		desc:    "connects to a peer",
		args:    "/connect [IP] (port)",
		example: []string{
			"/connect localhost",
			"/connect 192.160.1.24 8908",
		},
		callback: func(line string, args []string) {
			var ip = args[0] + ":"
			if args[1] != "" {
				ip += args[1]
			} else {
				ip += strconv.Itoa(common.DefaultPort)
			}
			logger.Println("Connecting to: " + ip)
			createConnFunc(ip)
		},
	},
	Command{
		regex:   regexp.MustCompile(`\/chats`),
		command: "chats",
		desc:    "lists chats",
		args:    "/chats",
		example: []string{
			"/chats",
		},
		callback: func(line string, args []string) {
			logger.Println("--- CHATS ---")
			for name, chat := range chatMap {
				if len(chat) > 0 {
					var printStr string
					printStr += "• " + name
					if len(chat) > 1 {
						printStr += " " + styles["group"]("(group)")
					}
					logger.Println(printStr)
				}
			}
		},
	},
	Command{
		regex:   regexp.MustCompile(`\/current`),
		command: "current",
		desc:    "displays current chat",
		args:    "/current",
		example: []string{
			"/current",
		},
		callback: func(line string, args []string) {
			var printString string
			printString += "Current Chat: " + currentChat + "\n"
			if currentChat != "" {
				printString += "Users: " + strings.Join(chatMap[currentChat], ", ")
			}
			logger.Println(printString)
		},
	},
	Command{
		regex:   regexp.MustCompile(`\/chat (.+)`),
		command: "chat",
		desc:    "switches to the given chat",
		args:    "/chat [chat name]",
		example: []string{
			"/chat slaidan_lt",
			"/chat leijurv",
		},
		callback: func(line string, args []string) {
			chat := args[0]
			hasChat := false
			for name := range chatMap {
				if name == chat {
					hasChat = true
					logger.Println("Connecting to chat", chat)
					currentChat = name
				}
			}
			if !hasChat {
				logger.Println("Error: Missing Chat")
			}
		},
	},
	Command{
		regex:   regexp.MustCompile(`\/group ([^ ]+) (.+)`),
		command: "group",
		desc:    "creates a group",
		args:    "/group [name] [users,here]",
		example: []string{
			"/chat slaidan_lt, leijurv",
			"/chat leijurv",
		},
		callback: func(line string, args []string) {
			groupName := strings.TrimSpace(args[0])
			usersArr := strings.Split(args[1], ",")
			for i := range usersArr {
				usersArr[i] = strings.TrimSpace(usersArr[i])
			}
			for name, arr := range chatMap {
				if name == groupName {
					logger.Println("Error: Chat Already Exists by That Name")
					return
				}
				if reflect.DeepEqual(usersArr, arr) {
					logger.Println("Error: Chat With the Same Users Already Exists")
					return
				}
			}
			AddGroup(groupName, usersArr)
			currentChat = groupName
		},
	},
}
