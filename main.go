package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	}
	app.UsageText = "chat [arguments...]"
	app.Action = runApp
	app.Run(os.Args)

}

func runApp(c *cli.Context) {
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	go printAll(outputChannel)
	go listen(c.Int("port"))
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1]
			if strings.Contains(text, "connect ") {
				ip := strings.Split(text, "connect ")[1] + ":8080"
				fmt.Println("Connecting " + ip)
				createConnection(ip)
			} else {
				fmt.Println("Sending")
				broadcastMessage(text)
			}
		}
	}()
	// Exit capture
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		err := saveConfig()
		if err != nil {
			panic(err)
		}
		done <- true
	}()
	<-done
	fmt.Println("Safe Exited")
}
