package cui

import (
	"github.com/billybobjoeaglt/chatlab/utils"
	"github.com/gizak/termui"
)

var windowMode int     // 0 = main window; 1 = commands; 2 = startup settings;
var selectedH int = 1  // 0 = chatList; 1 = text
var chatTextOffset int // 0 is latest message

// Show what is selected
func goLeft() {
	if selectedH == 1 {
		chatText.BorderFg = termui.ColorDefault
		chatList.BorderFg = termui.ColorYellow
		selectedH = 0
		reloadScreen()
	}
}

// Show what is selected
func goRight() {
	if selectedH == 0 {
		chatText.BorderFg = termui.ColorYellow
		chatList.BorderFg = termui.ColorDefault
		selectedH = 1
		reloadScreen()
	}
}

// If the list of chats is selected, go down a chat
// If the chat text is selected, scroll down to the latest message
func goDown() {
	if selectedH == 0 {
		index := utils.IndexOfStr(chatMapKeys, currentChat)

		// If it is possible to go down a chat
		if index != -1 && index < len(chatMapKeys)-1 {
			currentChat = chatMapKeys[index+1]

			// Load history with latest message a the bottom
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
		// If possible, go down a message
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

// If the list of chats is selected, go up a chat
// If the chat text is selected, scroll up to the oldest message
func goUp() {
	if selectedH == 0 {
		index := utils.IndexOfStr(chatMapKeys, currentChat)

		// If it is possible to go up a chat
		if index > 0 {
			currentChat = chatMapKeys[index-1]

			// Load history with latest message a the bottom
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
		// If possible, go up a message
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
