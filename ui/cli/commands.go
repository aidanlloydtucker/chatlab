package cli

import "regexp"

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
				ip += "8080"
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
			for _, chat := range chatList {
				logger.Println(chat)
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
			var chat string = "UNDEFINED"
			if len(chatList)-1 >= currentChat {
				chat = chatList[currentChat]
			}
			logger.Println("Current Chat:", chat)
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
			for i, val := range chatList {
				if val == chat {
					hasChat = true
					logger.Println("Connecting to chat", chat)
					currentChat = i
				}
			}
			if !hasChat {
				logger.Println("Error: Missing Chat")
			}
		},
	},
}
