package common

import (
	"io"
	"net/http"
	"os"
)

type Message struct {
	Username  string
	Decrypted bool
	Err       error
	Message   string
	Fullname  string
	ToUsers   []string
	ChatName  string
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

var Done = make(chan bool, 1)

var DefaultPort = 21991

func DoesUserExist(username string) (bool, error) {
	resp, err := http.Get("https://keybase.io/" + username)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode <= 300, nil
}

var ProgramDir string

func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
