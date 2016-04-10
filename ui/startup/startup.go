package startup

import (
	"log"
	"os"
	"path/filepath"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
	"gopkg.in/readline.v1"
)

// Run CLI dialog that sets variables and loads files
func RunStartup() {
	rl, err := readline.NewEx(&readline.Config{
		UniqueEditLine:  true,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// Assign logger variable to a new logger tied to STDOUT
	logger := log.New(rl.Stdout(), "", 0)

	// Check for username
	ok, err := common.DoesUserExist(config.GetConfig().Username)
	if err != nil {
		ok = false
	}
	if config.GetConfig().Username == "" || !ok {
		logger.Println("It seems you are missing your username. Please type it in")
		rl.SetPrompt("Username: ")
		un, err := rl.Readline()
		if err != nil {
			panic(err)
		}
		config.GetConfig().Username = un
	}

	// Check for private key
	if _, err := os.Stat(config.GetConfig().PrivateKey); os.IsNotExist(err) {
		logger.Println("It seems you are missing your private key. Please type in the file location")
		rl.SetPrompt("Private Key: ")
		pk, err := rl.Readline()
		if err != nil {
			panic(err)
		}
		config.GetConfig().PrivateKey = pk
	}

	if !config.GetConfig().AnsweredStorePK {
		logger.Println("Would you like us to store your private key in our program directory?")
		rl.SetPrompt("(y/N): ")
		for {
			yn, err := rl.Readline()
			if err != nil {
				panic(err)
			}
			if yn == "y" || yn == "N" {
				if yn == "y" {
					pkfp := filepath.Join(common.ProgramDir, "private.key")
					err := common.CopyFile(config.GetConfig().PrivateKey, pkfp)
					if err != nil {
						logger.Println("Error:", err.Error())
						return
					}
					config.GetConfig().PrivateKey = pkfp
				}
				config.GetConfig().AnsweredStorePK = true
				break
			} else {
				logger.Println("Error: Input is not valid")
			}
		}
	}

	// Check for password
	if config.GetConfig().Password == "" {
		if config.GetConfig().ShouldSavePass {
			logger.Println("It seems you are missing your password for your private key. Please type it in")
		}

		pass, err := rl.ReadPassword("Password: ")
		if err != nil {
			panic(err)
		}
		if config.GetConfig().ShouldSavePass {
			logger.Println("Would you like us to autosave this password?")
			rl.SetPrompt("(y/N): ")
			for {
				yn, err := rl.Readline()
				if err != nil {
					panic(err)
				}
				if yn == "y" || yn == "N" {
					if yn == "y" {
						config.GetConfig().Password = string(pass)
					} else {
						config.GetConfig().ShouldSavePass = false
					}
					config.Password = string(pass)
					break
				} else {
					logger.Println("Error: Input is not valid")
				}
			}
		}
	}

	if config.GetConfig().FirstTime {
		config.GetConfig().FirstTime = false
	}
}
