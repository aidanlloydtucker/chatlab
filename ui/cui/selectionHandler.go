package cui

import (
	"github.com/billybobjoeaglt/chatlab/utils"
	"github.com/gizak/termui"
)

var windowMode int     // 0 = main window; 1 = commands; 2 = startup settings;
var selectedH int = 1  // 0 = chatList; 1 = text
var chatTextOffset int // 0 is latest message

func goLeft() {
	if selectedH == 1 {
		chatText.BorderFg = termui.ColorDefault
		chatList.BorderFg = termui.ColorYellow
		selectedH = 0
		reloadScreen()
	}
}

func goRight() {
	if selectedH == 0 {
		chatText.BorderFg = termui.ColorYellow
		chatList.BorderFg = termui.ColorDefault
		selectedH = 1
		reloadScreen()
	}
}

func goDown() {
	if selectedH == 0 {
		index := utils.IndexOfStr(chatMapKeys, currentChat)
		if index != -1 && index < len(chatMapKeys)-1 {
			currentChat = chatMapKeys[index+1]
			var str string

			iStart := len(chatMap[currentChat].History) - (chatText.InnerHeight() - 1)
			if iStart < 0 {
				iStart = 0
			}
			for _, msg := range chatMap[currentChat].History[iStart:] {
				str += formatMsg(*msg.Message) + "\n"
			}
			chatText.Text = str
			reloadChatList()
		}
	} else if selectedH == 1 {
		iStart := len(chatMap[currentChat].History) - (chatText.InnerHeight() - 1)
		if iStart < 0 {
			iStart = 0
		}
		actualStart := iStart - chatTextOffset
		if iStart > 0 && chatTextOffset > 0 {
			chatTextOffset--
			var str string
			for _, msg := range chatMap[currentChat].History[actualStart+1 : len(chatMap[currentChat].History)-chatTextOffset] {
				str += formatMsg(*msg.Message) + "\n"
			}
			chatText.Text = str
			reloadScreen()
		}
	}
}

func goUp() {
	if selectedH == 0 {
		index := utils.IndexOfStr(chatMapKeys, currentChat)
		if index > 0 {
			currentChat = chatMapKeys[index-1]
			var str string

			iStart := len(chatMap[currentChat].History) - (chatText.InnerHeight() - 1)
			if iStart < 0 {
				iStart = 0
			}
			for _, msg := range chatMap[currentChat].History[iStart:] {
				str += formatMsg(*msg.Message) + "\n"
			}
			chatText.Text = str
			reloadChatList()
		}
	} else if selectedH == 1 {
		iStart := len(chatMap[currentChat].History) - (chatText.InnerHeight())
		if iStart < 0 {
			iStart = 0
		}
		actualStart := iStart - chatTextOffset
		if actualStart >= 0 && iStart > 0 {
			chatTextOffset++
			var str string
			for _, msg := range chatMap[currentChat].History[actualStart : len(chatMap[currentChat].History)-chatTextOffset] {
				str += formatMsg(*msg.Message) + "\n"
			}
			chatText.Text = str
			reloadScreen()
		}
	}
}
