package logger

// If true, lots of things will print
var IsVerbose bool

type ConsoleMessage struct {
	Level   Level // -1 verbose; 0 info; 1 important; 2 warning; 3 error;
	Message string
	Error   error
}

type ChanMessage chan ConsoleMessage
type ChanChanMessage chan ChanMessage

func (cc ChanMessage) AddVerbose(message string) {
	cc <- ConsoleMessage{Level: VERBOSE, Message: message}
}
func (ccm ChanChanMessage) AddVerbose(message string) {
	cc := make(ChanMessage)
	ccm <- cc
	cc.AddVerbose(message)
	close(cc)
}

func (cc ChanMessage) AddError(err error, message string) {
	cc <- ConsoleMessage{Level: ERROR, Message: message, Error: err}
}
func (ccm ChanChanMessage) AddError(err error, message string) {
	cc := make(ChanMessage)
	ccm <- cc
	cc.AddError(err, message)
	close(cc)
}

func (cc ChanMessage) AddInfo(message string) {
	cc <- ConsoleMessage{Level: INFO, Message: message}
}
func (ccm ChanChanMessage) AddInfo(message string) {
	cc := make(ChanMessage)
	ccm <- cc
	cc.AddInfo(message)
	close(cc)
}

func (cc ChanMessage) AddPriority(message string) {
	cc <- ConsoleMessage{Level: PRIORITY, Message: message}
}
func (ccm ChanChanMessage) AddPriority(message string) {
	cc := make(ChanMessage)
	ccm <- cc
	cc.AddInfo(message)
	close(cc)
}

type Level int

const (
	VERBOSE Level = iota
	INFO
	PRIORITY
	WARNING
	ERROR
)

var ConsoleChan = make(ChanChanMessage)
