package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/billybobjoeaglt/chatlab/chat"
	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
	"github.com/billybobjoeaglt/chatlab/logger"
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
			Value: common.DefaultPort,
			Usage: "set port of client",
		},
		/*cli.BoolFlag{
			Name:  "nogui, n",
			Usage: "Disables GUI",
		},*/
		cli.BoolFlag{
			Name:  "gui, g",
			Usage: "Enables GUI",
		},
		cli.BoolFlag{
			Name:  "cli, c",
			Usage: "Enables CLI",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enables verbosity",
		},
	}
	app.UsageText = "chat [arguments...]"
	app.Action = runApp
	app.Run(os.Args)

}

func runApp(c *cli.Context) {
	var err error
	err = config.LoadConfig()
	if err != nil {
		panic(err)
	}
	go uiPrint(chat.GetMessageChannel())
	//go printAll(chat.GetOutputChannel())
	go chat.Listen(c.Int("port"))

	var ip string = "UNKNOWN"

	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ip = net.JoinHostPort(ipnet.IP.String(), strconv.Itoa(c.Int("port")))
			}
		}
	}

	fmt.Println("Broadcasting on: " + ip)

	if c.Bool("gui") {
		ui.NewGUI()
	} else if c.Bool("cli") {
		err = ui.NewCLI()
		if err != nil {
			panic(err)
		}
	}

	logger.Verbose = c.Bool("verbose")

	ui.SetSendMessage(func(msg common.Message) {
		go chat.BroadcastMessage(msg)
	})
	ui.SetCreateConn(func(ip string) {
		go chat.CreateConnection(ip)
	})

	// Exit capture
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		common.Done <- true
	}()
	<-common.Done
	err = config.SaveConfig()
	if err != nil {
		panic(err)
	}
	ui.Quit()
	fmt.Println("Safe Exited")
}

/*func printAll(stringChanChan <-chan chan string) {
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
}*/
func uiPrint(msgChan <-chan common.Message) {
	for {
		msg, ok := <-msgChan
		if ok {
			ui.AddMessage(msg)
		}
	}
}
