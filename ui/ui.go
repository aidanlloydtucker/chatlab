package ui

import (
	"log"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/logger"
	"github.com/billybobjoeaglt/chatlab/ui/cli"
	"github.com/billybobjoeaglt/chatlab/ui/cui"
)

// Defines the type of ui running.
// 0 is none; 1 is CLI; 2 is GUI;
var uiType int

// A STDOUT console for no UI
func RelayConsole(ccChan *logger.ChanChanMessage) {
	for {
		cc := <-*ccChan
		for {
			cm, ok := <-cc
			if ok {
				switch cm.Level {
				case logger.VERBOSE:
					if logger.IsVerbose {
						log.Println("VERBOSE:", cm.Message)
					}
				case logger.INFO:
					log.Println("INFO:", cm.Message)
				case logger.PRIORITY:
					log.Println("IMPORTANT:", cm.Message)
				case logger.WARNING:
					log.Println("WARNING:", cm.Message)
				case logger.ERROR:
					log.Println("ERROR:", cm.Message)
					log.Println(cm.Error.Error())
				}
			} else {
				break
			}
		}
	}
}

func NewRelayConsole() {
	go RelayConsole(&logger.ConsoleChan)
}

// Creates new CLI
func NewCLI() error {
	if uiType == 0 {
		go cli.StartCLI()
		go cli.CLIConsole(&logger.ConsoleChan)
		uiType = 1
	}
	return nil
}

// Creates new CUI
func NewCUI() error {
	if uiType == 0 {
		go cui.StartCUI()
		go cui.CUIConsole(&logger.ConsoleChan)
		uiType = 2
	}
	return nil
}

// Sets the chat function for sending message
func SetSendMessage(f common.SendMessageFunc) {
	switch uiType {
	case 1:
		cli.SetSendMessage(f)
	case 2:
		cui.SetSendMessage(f)
	}
}

// Sets the chat function for creating connection
func SetCreateConn(f common.CreateConnFunc) {
	switch uiType {
	case 1:
		cli.SetCreateConn(f)
	case 2:
		cui.SetCreateConn(f)
	}
}

// Quits UI
func Quit() {
	switch uiType {
	case 1:
		cli.QuitCLI()
	case 2:
		cui.QuitCUI()
	}
}

// Adds a message to UI
func AddMessage(msg common.Message) {
	switch uiType {
	case 1:
		cli.AddMessage(msg)
	case 2:
		cui.AddMessage(msg)
	}
}

// Adds a user to UI
func AddUser(user string) {
	switch uiType {
	case 1:
		cli.AddUser(user)
	case 2:
		cui.AddUser(user)
	}
}

// Removes a user from UI
func RemoveUser(user string) {
	switch uiType {
	case 1:
		cli.RemoveUser(user)
	case 2:
		cui.RemoveUser(user)
	}
}

// Adds a group to UI
func AddGroup(groupName string, users []string) {
	switch uiType {
	case 1:
		cli.AddGroup(groupName, users)
	case 2:
		cui.AddGroup(groupName, users)
	}
}
