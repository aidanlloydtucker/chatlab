package main

import "os"

var config *Config

type Config struct {
	username   string
	privateKey string
	passphrase string
}

func makeConfig() (*Config, error) {
	var newConfig = &Config{
		username:   "user",
		privateKey: "",
		passphrase: "",
	}
	return newConfig, nil
}

func loadConfig() error {
	configFile, err := os.Open("./config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()
	if _, err := configFile.Stat(); !os.IsNotExist(err) {
		newConfig, err := makeConfig()
		config = newConfig
		if err != nil {
			return err
		}
	}
	return nil
}
