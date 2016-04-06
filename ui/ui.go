package ui

import (
	"fmt"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/ui/cli"
)

// Defines the type of ui running.
// 0 is none; 1 is CLI; 2 is GUI;
var uiType int

// Creates new GUI
func NewGUI() {
	if uiType == 0 {
		//gui.StartGUI()
		uiType = 2
	}
}

// Creates new CLI
func NewCLI() error {
	if uiType == 0 {
		go cli.StartCLI()
		uiType = 1
	}
	return nil
}

// Sets the chat function for sending message
func SetSendMessage(f common.SendMessageFunc) {
	switch uiType {
	case 1:
		cli.SetSendMessage(f)
	case 2:
		//gui.SetSendMessage(f)
	}
}

// Sets the chat function for creating connection
func SetCreateConn(f common.CreateConnFunc) {
	switch uiType {
	case 1:
		cli.SetCreateConn(f)
	case 2:
		//gui.SetSendMessage(f)
	}
}

// Quits UI
func Quit() {
	switch uiType {
	case 1:
		cli.QuitCLI()
	case 2:
		//gui.QuitGUI()
	}
}

// Adds a message to UI
func AddMessage(msg common.Message) {
	switch uiType {
	case 1:
		cli.AddMessage(msg)
	case 2:
		//gui.AddMessage(user, message)
		fmt.Println(msg.Username + ": " + msg.Message)
	}
}

// Adds a user to UI
func AddUser(user string) {
	switch uiType {
	case 1:
		cli.AddUser(user)
	case 2:
		//gui.AddUser(user)
	}
}

// Removes a user from UI
func RemoveUser(user string) {
	switch uiType {
	case 1:
		cli.RemoveUser(user)
	}
}

// Adds a group to UI
func AddGroup(groupName string, users []string) {
	switch uiType {
	case 1:
		cli.AddGroup(groupName, users)
	}
}
