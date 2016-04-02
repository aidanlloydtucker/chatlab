package gui

import (
	"unsafe"

	"github.com/billybobjoeaglt/chatlab/ui/common"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var chatBoxBuf *gtk.TextBuffer
var chatBoxIter gtk.TextIter
var chatSelectList *gtk.ListStore
var chatSelectIter gtk.TreeIter
var sendMessageCB common.SendMessageFunc

func StartGUI() {
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("ChatLab GUI")

	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	}, "foo")

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
	msgField.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		key := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		if key.Keyval == 65293 {
			go sendMessageFromGUI(msgField)
		}
	})
	hbox.Add(msgField)

	sendButton := gtk.NewButtonWithLabel("Send")
	sendButton.Clicked(func() {
		go sendMessageFromGUI(msgField)
	})
	hbox.Add(sendButton)

	vbox.Add(swin)
	vbox.Add(hbox)

	hpane.Add2(vbox)

	window.Add(hpane)

	window.SetSizeRequest(1080, 720)
	window.ShowAll()
	gtk.Main()

}
func QuitGUI() {
	gtk.MainQuit()
}

func sendMessageFromGUI(msgField *gtk.Entry) {
	if msgField.GetTextLength() == 0 {
		return
	}
	if sendMessageCB != nil {
		sendMessageCB("bob", msgField.GetText())
	}
	AddMessage("bob", msgField.GetText())
	msgField.SetText("")
}

func SetSendMessage(f common.SendMessageFunc) {
	sendMessageCB = f
}

func AddMessage(user string, message string) {
	gdk.ThreadsEnter()
	chatBoxBuf.Insert(&chatBoxIter, user+": "+message+"\n")
	gdk.ThreadsLeave()
}

func AddUser(user string) {
	gdk.ThreadsEnter()
	chatSelectList.Append(&chatSelectIter)
	chatSelectList.Set(&chatSelectIter,
		0, user)
	gdk.ThreadsLeave()
}
