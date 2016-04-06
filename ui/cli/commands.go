package cli

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/billybobjoeaglt/chatlab/common"
)

type CommandCallback func(line string, args []string)

type Command struct {
	Regex    *regexp.Regexp
	Command  string
	Desc     string
	Args     string
	Example  []string
	Callback CommandCallback
}

var commandArr = []Command{
	Command{
		Regex:   regexp.MustCompile(`\/connect ([^ ]+) ?(.*)`),
		Command: "connect",
		Desc:    "connects to a peer",
		Args:    "/connect [IP] (port)",
		Example: []string{
			"/connect localhost",
			"/connect 192.160.1.24 8908",
		},
		Callback: func(line string, args []string) {
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
		Regex:   regexp.MustCompile(`\/chats`),
		Command: "chats",
		Desc:    "lists chats",
		Args:    "/chats",
		Example: []string{
			"/chats",
		},
		Callback: func(line string, args []string) {
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
		Regex:   regexp.MustCompile(`\/current`),
		Command: "current",
		Desc:    "displays current chat",
		Args:    "/current",
		Example: []string{
			"/current",
		},
		Callback: func(line string, args []string) {
			var printString string
			printString += "Current Chat: " + currentChat + "\n"
			if currentChat != "" {
				printString += "Users: " + strings.Join(chatMap[currentChat], ", ")
			}
			logger.Println(printString)
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/chat (.+)`),
		Command: "chat",
		Desc:    "switches to the given chat",
		Args:    "/chat [chat name]",
		Example: []string{
			"/chat slaidan_lt",
			"/chat leijurv",
		},
		Callback: func(line string, args []string) {
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
		Regex:   regexp.MustCompile(`\/group ([^ ]+) (.+)`),
		Command: "group",
		Desc:    "creates a group",
		Args:    "/group [name] [users,here]",
		Example: []string{
			"/chat slaidan_lt, leijurv",
			"/chat leijurv",
		},
		Callback: func(line string, args []string) {
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
