package ui

import (
	"fmt"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/ui/cli"
	"github.com/billybobjoeaglt/chatlab/ui/gui"
)

var uiType int // 0 is none; 1 is CLI; 2 is GUI;

func NewGUI() {
	if uiType == 0 {
		gui.StartGUI()
		uiType = 2
	}
}

func NewCLI() error {
	if uiType == 0 {
		go cli.StartCLI()
		uiType = 1
	}
	return nil
}

func SetSendMessage(f common.SendMessageFunc) {
	switch uiType {
	case 1:
		cli.SetSendMessage(f)
	case 2:
		gui.SetSendMessage(f)
	}
}

func SetCreateConn(f common.CreateConnFunc) {
	switch uiType {
	case 1:
		cli.SetCreateConn(f)
	case 2:
		//gui.SetSendMessage(f)
	}
}

func Quit() {
	switch uiType {
	case 1:
		cli.QuitCLI()
	case 2:
		gui.QuitGUI()
	}
}

/*func AddMessageChan(stringChanChan <-chan chan string) {
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
}*/

func AddMessage(msg common.Message) {
	switch uiType {
	case 1:
		cli.AddMessage(msg)
	case 2:
		//gui.AddMessage(user, message)
		fmt.Println(msg.Username + ": " + msg.Message)
	}
}

func AddUser(user string) {
	switch uiType {
	case 1:
		cli.AddUser(user)
	case 2:
		gui.AddUser(user)
	}
}
func RemoveUser(user string) {
	switch uiType {
	case 1:
		cli.RemoveUser(user)
	}
}
