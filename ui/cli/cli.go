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
)

var logger *log.Logger

var sendMsgFunc common.SendMessageFunc
var createConnFunc common.CreateConnFunc
var username *string

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
	if strings.HasPrefix(line, "/") {
		fmt.Println("COMMAND:", line)
		switch {
		case strings.Contains(line, "/connect"):
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
		}
	} else {
		if sendMsgFunc != nil {
			msg := common.NewMessage()
			msg.Username = *username
			msg.Message = line
			msg.ToUsers = append(msg.ToUsers, "slaidan_lt")
			msg.ToUsers = append(msg.ToUsers, "leijurv")
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
	logger.Println("Quitting")
}

func AddMessage(msg common.Message) {
	logger.Println(msg.Username + ": " + msg.Message)
}

func AddUser(user string) {
	logger.Println("New User:", user)
}

/*go func() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		if strings.Contains(text, "connect ") {
			ip := strings.Split(text, "connect ")[1] + ":8080"
			fmt.Println("Connecting " + ip)
			chat.CreateConnection(ip)
		} else {
			fmt.Println("Sending")
			chat.BroadcastMessage(text)
		}
	}
}()*/
