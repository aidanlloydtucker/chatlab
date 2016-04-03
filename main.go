package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/billybobjoeaglt/chatlab/chat"
	"github.com/billybobjoeaglt/chatlab/config"
	"github.com/billybobjoeaglt/chatlab/ui"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "ChatLab"
	app.Usage = "A P2P Encrypted Chat App"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Aidan Lloyd-Tucker",
		},
		cli.Author{
			Name: "Leif",
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8080,
			Usage: "set port of client",
		},
		cli.BoolFlag{
			Name:  "nogui, n",
			Usage: "Disables GUI",
		},
	}
	app.UsageText = "chat [arguments...]"
	app.Action = runApp
	app.Run(os.Args)

}

func runApp(c *cli.Context) {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	go printAll(chat.GetOutputChannel())
	go chat.Listen(c.Int("port"))
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1]
			if strings.Contains(text, "connect ") {
				ip := strings.Split(text, "connect ")[1] + ":8080"
				fmt.Println("Connecting " + ip)
				chat.CreateConnection(ip)
			} else {
				fmt.Println("Sending")
				chat.BroadcastMessage(text)
			}
		}
	}()

	if !c.Bool("nogui") {
		ui.NewGUI()
	}
	ui.SetSendMessage(func(user string, message string) {
		go chat.BroadcastMessage(message)
	})
	//go addMessageUI(chat.GetOutputChannel())

	// Exit capture
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		err := config.SaveConfig()
		if err != nil {
			panic(err)
		}
		ui.Quit()
		done <- true
	}()
	<-done
	fmt.Println("Safe Exited")
}

func printAll(stringChanChan <-chan chan string) {
	for {
		strChan := <-stringChanChan
		for {
			str, ok := <-strChan
			if ok {
				fmt.Printf(str)
			} else {
				break
			}
		}
		fmt.Println()
	}
}
