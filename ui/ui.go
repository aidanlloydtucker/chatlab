package ui

import (
	"fmt"

	"github.com/billybobjoeaglt/chatlab/ui/common"
	"github.com/billybobjoeaglt/chatlab/ui/gui"
)

var hasUI bool = false

func NewGUI() {
	if !hasUI {
		gui.StartGUI()
		hasUI = true
	}
}

func SetSendMessage(f common.SendMessageFunc) {
	if hasUI {
		gui.SetSendMessage(f)
	}
}

func Quit() {
	if hasUI {
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

func AddMessage(user string, message string) {
	if hasUI {
		fmt.Println(user + ": " + message)
	}
}

func AddUser(user string) {
	if hasUI {
		gui.AddUser(user)
	}
}
