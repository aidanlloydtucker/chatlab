package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Decrypted Message struct
type Message struct {
	Username  string
	Decrypted bool
	Err       error
	Message   string
	Fullname  string
	ToUsers   []string
	ChatName  string
}

// Create a new message with defaults
func NewMessage() *Message {
	return &Message{
		Username:  "",
		Decrypted: true,
		Err:       nil,
		Message:   "",
		Fullname:  "",
	}
}

// Send Message Function from UI type
type SendMessageFunc func(Message)

// Create Connection Function from UI type
type CreateConnFunc func(string)

// Make this true if you want to quit
var Done = make(chan bool, 1)

var DefaultPort = 21991

// Check if user exists
// TODO: Protect against usernames like "/" or " ", that could just go to keybase.io
func DoesUserExist(username string) (bool, error) {
	if username == "" {
		return false, nil
	}
	resp, err := http.Get("https://keybase.io/" + username)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode <= 300, nil
}

// Where the default program directory is
var ProgramDir string

// NOTE: I *NEVER* copy from stack overflow
// XXX: Link for me later: http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
