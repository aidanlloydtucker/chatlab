package cli

import (
	"fmt"
	"io"
	"log"
	"strings"
	"syscall"

	"gopkg.in/readline.v1"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
	"github.com/ttacon/chalk"
)

var logger *log.Logger

var sendMsgFunc common.SendMessageFunc
var createConnFunc common.CreateConnFunc
var username *string
var chatList []string
var currentChat int

var styles = map[string]func(string) string{
	"username":     chalk.Blue.NewStyle().WithTextStyle(chalk.Bold).Style,
	"notification": chalk.White.NewStyle().WithTextStyle(chalk.Underline).Style,
}

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

func StartCLI() {

	username = &config.GetConfig().Username

	rl, err := readline.NewEx(&readline.Config{
		UniqueEditLine:  true,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	logger = log.New(rl.Stderr(), "", 0)

	rl.SetPrompt("> ")
	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				syscall.Exit(0)
			} else if err == io.EOF {
				syscall.Exit(0)
			}
			break
		}
		lineHandler(line)
	}
}

func lineHandler(line string) {
	if line == "" {
		return
	}
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "/") {
		switch {
		case strings.HasPrefix(line, "/connect "):
			var ip string
			args := strings.Split(line, " ")[1:]
			switch len(args) {
			case 0:
				logger.Println("Error: Missing IP")
				return
			case 1:
				ip = args[0] + ":8080"
			case 2:
				ip = args[0] + ":" + args[1]
			default:
				logger.Println("Error: Too Many Args")
				return
			}
			logger.Println("Connecting to: " + ip)
			createConnFunc(ip)
		case line == "/chats":
			logger.Println("--- CHATS ---")
			for _, chat := range chatList {
				logger.Println(chat)
			}
		case line == "/current":
			var chat string = "UNDEFINED"
			if len(chatList)-1 >= currentChat {
				chat = chatList[currentChat]
			}
			logger.Println("Current Chat:", chat)
		case strings.HasPrefix(line, "/chat "):
			args := strings.Split(line, " ")[1:]
			switch len(args) {
			case 0:
				logger.Println("Error: Missing Chat")
				return
			case 1:
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
				return
			default:
				logger.Println("Error: Too Many Args")
				return
			}
		}
	} else {
		if sendMsgFunc != nil {
			msg := common.NewMessage()
			msg.Username = *username
			msg.Message = line
			msg.ToUsers = append(msg.ToUsers, chatList[currentChat])
			AddMessage(*msg)
			sendMsgFunc(*msg)
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
	logger.Println(styles["username"](msg.Username+":"), msg.Message)
}

func AddUser(user string) {
	chatList = append(chatList, user)
	logger.Println(styles["notification"]("New User: " + user))
}
