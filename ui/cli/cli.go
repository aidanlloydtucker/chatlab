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

var logger *log.Logger

var sendMsgFunc common.SendMessageFunc
var createConnFunc common.CreateConnFunc
var username *string
var chatList []string
var currentChat int

var styles = map[string]func(string) string{
	"username":     chalk.Blue.NewStyle().WithTextStyle(chalk.Bold).Style,
	"notification": chalk.White.NewStyle().WithTextStyle(chalk.Underline).Style,
	"error":        chalk.Red.NewStyle().WithTextStyle(chalk.Underline).Style,
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
				common.Done <- true
			} else if err == io.EOF {
				common.Done <- true
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
		msg := common.NewMessage()
		msg.Username = *username
		msg.Message = line
		AddMessage(*msg)

		noMatches := true
		for _, cmd := range commandArr {
			if cmd.regex.MatchString(line) {
				cmd.callback(line, cmd.regex.FindStringSubmatch(line)[1:])
				noMatches = false
			}
		}
		if noMatches {
			logger.Println(styles["error"]("Error: Unknown Command"))
		}
	} else {
		if sendMsgFunc != nil && len(chatList)-1 >= currentChat {
			msg := common.NewMessage()
			msg.Username = *username
			msg.Message = line
			msg.ToUsers = append(msg.ToUsers, chatList[currentChat])
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
	logger.Println(styles["username"](msg.Username+":"), msg.Message)
}

func RemoveUser(user string) {
	for i, val := range chatList {
		if val == user {
			chatList = chatList[:i+copy(chatList[i:], chatList[i+1:])]
			break
		}
	}
}

func AddUser(user string) {
	chatList = append(chatList, user)
	logger.Println(styles["notification"]("New User: " + user))
}
