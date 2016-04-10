package cui

import (
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/billybobjoeaglt/chatlab/common"
)

type CommandCallback func(line string, args []string)

type Command struct {
	Regex    *regexp.Regexp
	Command  string
	Desc     string
	Args     string
	Example  []string
	Callback CommandCallback
}

var commandArr = []Command{
	Command{
		Regex:   regexp.MustCompile(`\/connect ([^ ]+) ?(.*)`),
		Command: "connect",
		Desc:    "connects to a peer",
		Args:    "/connect [IP] (port)",
		Example: []string{
			"/connect localhost",
			"/connect 192.160.1.24 8908",
		},
		Callback: func(line string, args []string) {
			var port string
			if args[1] != "" {
				port = args[1]
			} else {
				port = strconv.Itoa(common.DefaultPort)
			}
			createConnFunc(net.JoinHostPort(args[0], port))
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/group ([^ ]+) (.+)`),
		Command: "group",
		Desc:    "creates a group",
		Args:    "/group [name] [users,here]",
		Example: []string{
			"/chat slaidan_lt, leijurv",
			"/chat leijurv",
		},
		Callback: func(line string, args []string) {
			groupName := strings.TrimSpace(args[0])
			usersArr := strings.Split(args[1], ",")
			for i := range usersArr {
				usersArr[i] = strings.TrimSpace(usersArr[i])
			}
			for name, arr := range chatMap {
				if name == groupName {
					printError("Chat Already Exists by That Name")
					return
				}
				if reflect.DeepEqual(usersArr, arr) {
					printError("Chat With the Same Users Already Exists")
					return
				}
			}
			AddGroup(groupName, usersArr)
			currentChat = groupName
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/user (.+)`),
		Command: "user",
		Desc:    "adds a user",
		Args:    "/user [username]",
		Example: []string{
			"/user slaidan_lt",
			"/user leijurv",
		},
		Callback: func(line string, args []string) {
			username := strings.TrimSpace(args[0])
			userArr := []string{username}
			for name, chat := range chatMap {
				if name == username {
					printError("User Already Exists by That Name")
					return
				}
				if reflect.DeepEqual(userArr, chat.Users) {
					printError("Chat With the User Already Exists")
					return
				}
			}
			ok, err := common.DoesUserExist(username)
			if err != nil {
				printError(err.Error())
				return
			}
			if !ok {
				printError("User Does Not Exist")
				return
			}
			AddUser(username)
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/quit`),
		Command: "quit",
		Desc:    "quits chatlab",
		Args:    "/quit",
		Example: []string{
			"/quit",
		},
		Callback: func(line string, args []string) {
			common.Done <- true
		},
	},
}
