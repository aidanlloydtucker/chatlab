package ui

import (
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

func AddMessage(user string, message string) {
	if hasUI {
		gui.AddMessage(user, message)
	}
}

func AddUser(user string) {
	if hasUI {
		gui.AddUser(user)
	}
}
