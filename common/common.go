package common

type Message struct {
	Username  string
	Decrypted bool
	Err       error
	Message   string
	Fullname  string
	ToUsers   []string
}

func NewMessage() *Message {
	return &Message{
		Username:  "",
		Decrypted: true,
		Err:       nil,
		Message:   "",
		Fullname:  "",
	}
}

type SendMessageFunc func(Message)
type CreateConnFunc func(string)
