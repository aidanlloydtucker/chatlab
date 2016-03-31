package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var config *Config

type Config struct {
	Username   string
	PrivateKey string
	Passphrase string
}

func makeConfig() (*Config, error) {
	var newConfig = &Config{
		Username:   "user",
		PrivateKey: "./key.key",
		Passphrase: "./pass.key",
	}
	return newConfig, nil
}

func loadConfig() error {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		var newConfig *Config
		newConfig, err = makeConfig()
		if err != nil {
			return err
		}
		config = newConfig
	} else {
		configFile, err := os.Open("./config.json")
		if err != nil {
			return err
		}
		defer configFile.Close()
		byteArrFile, err := ioutil.ReadAll(configFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(byteArrFile, &config)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveConfig() error {
	file, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./config.json", file, 0777)
	return err
}
