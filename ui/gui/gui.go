package gui

import (
	"fmt"

	"github.com/billybobjoeaglt/chatlab/common"
)

/*var chatBoxBuf *gtk.TextBuffer
var chatBoxIter gtk.TextIter
var chatSelectList *gtk.ListStore
var chatSelectIter gtk.TreeIter*/
var sendMessageCB common.SendMessageFunc

func StartGUI() {
	/*gdk.ThreadsInit()
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("ChatLab GUI")

	gdk.ThreadsEnter()
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	}, "foo")
	gdk.ThreadsLeave()

	hpane := gtk.NewHPaned()
	hpane.SetPosition(200)

	chatSelectList = gtk.NewListStore(glib.G_TYPE_STRING)

	chatSelect := gtk.NewTreeView()
	chatSelect.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Name", gtk.NewCellRendererText(), "text", 0))
	chatSelect.SetModel(chatSelectList)

	hpane.Add1(chatSelect)

	vbox := gtk.NewVBox(false, 5)

	// Scroll Window
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)

	// Chat Box
	chatBox := gtk.NewTextView()
	chatBox.SetEditable(false)
	chatBox.SetCursorVisible(false)
	chatBox.SetSizeRequest(-1, 600)

	chatBoxBuf = chatBox.GetBuffer()

	var chatBoxStartIter gtk.TextIter

	chatBoxBuf.GetStartIter(&chatBoxStartIter)
	chatBoxBuf.GetEndIter(&chatBoxIter)
	swin.Add(chatBox)

	hbox := gtk.NewHBox(false, 5)

	msgField := gtk.NewEntry()
	gdk.ThreadsEnter()
	msgField.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		key := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		if key.Keyval == 65293 {
			go sendMessage(msgField)
		}
	})
	hbox.Add(msgField)

	sendButton := gtk.NewButtonWithLabel("Send")
	sendButton.Clicked(func() {
		go sendMessage(msgField)
	})
	gdk.ThreadsLeave()
	hbox.Add(sendButton)

	vbox.Add(swin)
	vbox.Add(hbox)

	hpane.Add2(vbox)

	window.Add(hpane)

	window.SetSizeRequest(1080, 720)
	window.ShowAll()

	gdk.ThreadsEnter()
	go gtk.Main()
	gdk.ThreadsLeave()*/

}
func QuitGUI() {
	//gtk.MainQuit()
}

/*func sendMessage(msgField *gtk.Entry) {
	if msgField.GetTextLength() > 0 {
		/*if sendMessageCB != nil {
			sendMessageCB("bob", msgField.GetText())
		}*
		AddMessage(config.GetConfig().Username + ": " + msgField.GetText())
		msgField.SetText("")
	}
}*/

func SetSendMessage(f common.SendMessageFunc) {
	sendMessageCB = f
}

func AddMessage(message string) {
	fmt.Println("adding")
	//gdk.ThreadsEnter()
	//chatBoxBuf.Insert(&chatBoxIter, message+"\n")
	//gdk.ThreadsLeave()
}

func AddUser(user string) {
	/*gdk.ThreadsEnter()
	chatSelectList.Append(&chatSelectIter)
	chatSelectList.Set(&chatSelectIter,
		0, user)
	gdk.ThreadsLeave()*/
}
