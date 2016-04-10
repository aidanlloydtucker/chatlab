package logger

// If true, lots of things will print
var IsVerbose bool

// Format for messages that will be printed on the console
type ConsoleMessage struct {
	Level   Level // -1 verbose; 0 info; 1 important; 2 warning; 3 error;
	Message string
	Error   error
}

// Channel of ConsoleMessages
type ChanMessage chan ConsoleMessage

// Channel of Channels of ConsoleMessages
// Used to maintain order
// Ex: A create connection thread is started, all messages related to the
// creation of that connection will go in its own channel, which will then
// go in a channel full of other channels
type ChanChanMessage chan ChanMessage

// Add a verbose message to a sub-channel
func (cc ChanMessage) AddVerbose(message string) {
	cc <- ConsoleMessage{Level: VERBOSE, Message: message}
}

// Add a verbose message to the main channel
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

// Logger level enum
type Level int

const (
	VERBOSE Level = iota
	INFO
	PRIORITY
	WARNING
	ERROR
)

// The main channel for logging
var ConsoleChan = make(ChanChanMessage)
