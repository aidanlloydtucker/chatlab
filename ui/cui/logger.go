package cui

import (
	"strings"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/logger"
)

func CUIConsole(ccChan *logger.ChanChanMessage) {
	for {
		cc := <-*ccChan
		for {
			cm, ok := <-cc
			if ok {
				msg := common.NewMessage()
				msg.ChatName = "logs"
				var msgText string
				switch cm.Level {
				case logger.VERBOSE:
					msg.Username = "VERBOSE"
					msgText = cm.Message
				case logger.INFO:
					msg.Username = "INFO"
					msgText = "[" + cm.Message + "](fg-white,fg-underline)"
				case logger.PRIORITY:
					msg.Username = "INFO"
					msgText = "[" + cm.Message + "](fg-white,fg-underline)"
					printInfo(cm.Message)
				case logger.WARNING:
					msg.Username = "WARNING"
					msgText = "[" + cm.Message + "](fg-yellow,fg-underline)"
					printWarning(cm.Message)
				case logger.ERROR:
					msg.Username = "ERROR"
					msgText = "[" + cm.Message + "](fg-red,fg-underline)"
					printError(cm.Message)
				}
				msg.Message = msgText
				if logger.IsVerbose || cm.Level != logger.VERBOSE {
					tmp := chatMap["logs"]
					tmp.History = append(tmp.History, Message{Message: msg, Read: msg.ChatName == currentChat})
					chatMap["logs"] = tmp

					if currentChat == "logs" || currentChat == "" {
						chatTextInit += formatMsg(*msg) + "\n"
					}
					reloadChatList()
				}
			} else {
				break
			}
		}
	}
}

func printInfo(str string) {
	printLn("[INFO: " + str + "](fg-white,fg-underline)")
}

func printWarning(str string) {
	printLn("[WARNING: " + str + "](fg-yellow,fg-underline)")
}

func printError(str string) {
	printLn("[ERROR: " + str + "](fg-red,fg-underline)")
}

func formatMsg(msg common.Message) string {
	var strMsg string
	strMsg += "[" + msg.Username + ":](fg-bold) " + msg.Message
	return strMsg
}

func printLn(str string) {
	if !uiMade {
		return
	}
	chatText.Text += str + "\n"
	chatTextArr := strings.Split(chatText.Text, "\n")
	if len(chatTextArr) > chatText.InnerHeight() {
		chatTextArr = append(chatTextArr[:0], chatTextArr[1:]...)
		chatText.Text = strings.Join(chatTextArr, "\n")
	}
	reloadScreen()
}

// Adds a command to logger
func AddCommand(msg common.Message) {
	printLn(msg.Username + ": " + msg.Message)
}
